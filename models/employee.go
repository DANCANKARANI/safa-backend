package models

import (
	"errors"
	"log"
	"time"

	"github.com/dancankarani/safa/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func (e *Employee) TableName() string {
	return "employees"
}

func NewEmployee(firstName, lastName, position, phoneNumber, email, address string) *Employee {
	return &Employee{
		ID:         uuid.New(),
		FirstName:  firstName,
		LastName:   lastName,
		Position:   position,
		PhoneNumber: phoneNumber,
		Email:      email,
	}
}

func (e *Employee) SetLoginCredentials() {
	e.Password,_ = services.GenerateFormattedPassword()
	e.CanLogin = true
}

func (e *Employee) SetRole(role string) {
	if role == "" {
		e.Role = "employee" // Default role
	} else {
		e.Role = role
	}
}



func (e *Employee) SetDateJoined(dateJoined string) {
	if dateJoined == "" {
		e.DateJoined = time.Now() // Default to current time if not provided
	} else {
		// Parse the date string and set DateJoined
		parsedDate, err := time.Parse("2006-01-02", dateJoined)
		if err == nil {
			e.DateJoined = parsedDate
		} else {
			e.DateJoined = time.Now() // Fallback to current time on error
		}
	}
}

func (e *Employee) BeforeCreate() {
	e.ID = uuid.New() // Ensure ID is set before creating
	e.CreatedAt = time.Now()
	e.UpdatedAt = time.Now()
}

func (e *Employee) BeforeUpdate() {
	e.UpdatedAt = time.Now() // Update timestamp before updating
}

func CreateEmployee(c *fiber.Ctx, e *Employee)( *Employee, error) {
	db.AutoMigrate(&e)
	e.ID = uuid.New()
	isValidPhone := services.ValidatePhoneNumber(e.PhoneNumber)
	if !isValidPhone {
		return nil, errors.New("invalid phone number")
	}

	isValidEmail := services.ValidateEmail(e.Email)
	if !isValidEmail {
		return nil, errors.New("invalid email")
	}

	password,_ := services.GenerateFormattedPassword()
	hashed_password, err := services.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	
	e.Password = hashed_password

	if err := db.Create(e).Error; err != nil {
		return nil, errors.New("failed to create employee")
	}
	log.Println("Pasword:", password)
	return e, nil
}

func GetEmployeeByID(c *fiber.Ctx, id uuid.UUID) (*Employee, error) {
	var employee Employee
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		return nil, errors.New("Employee not found")
	}
	return &employee, nil
}

func GetAllEmployees(c *fiber.Ctx) ([]Employee, error) {
	var employees []Employee
	if err := db.Find(&employees).Error; err != nil {
		return nil, errors.New("failed to retrieve employees")
	}
	return employees, nil
}

func UpdateEmployee(c *fiber.Ctx, id uuid.UUID, updatedData *Employee) (*Employee, error) {
	var employee Employee
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		return nil, errors.New("Employee not found")
	}
	if updatedData.PhoneNumber != "" {
		isValidPhone := services.ValidatePhoneNumber(updatedData.PhoneNumber)
		if !isValidPhone {
			return nil, errors.New("invalid phone number")
		}
	}

	if updatedData.Email != "" {
		isValidEmail := services.ValidateEmail(updatedData.Email)
		if !isValidEmail {
			return nil, errors.New("invalid email")
		}
	}

	//hash the password if it is provided
	hashedPassword, err := services.HashPassword(updatedData.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}
	updatedData.Password = hashedPassword
	employee.FirstName = updatedData.FirstName
	employee.LastName = updatedData.LastName
	employee.Position = updatedData.Position
	employee.PhoneNumber = updatedData.PhoneNumber
	employee.Email = updatedData.Email
	employee.Salary = updatedData.Salary
	employee.UpdatedAt = time.Now()
	if err := db.Save(&employee).Error; err != nil {
		return nil, errors.New("failed to update employee")
	}
	return &employee, nil
}

func DeleteEmployee(c *fiber.Ctx, id uuid.UUID) error {
	var employee Employee
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		return errors.New("employee not found")
	}
	if err := db.Where("id = ?", id).Delete(&employee).Error; err != nil {
		log.Println("failed to delete employee:", err.Error())
		return errors.New(err.Error())
	}
	return nil
}

//Get Employee Latest payments and All advances in the last one month
type ResEmployee struct {
	Employees Employee
	TotalPayments float64
	TotalAdvances float64
}
func GetEmployeePaymentsAndAdvances() ([]ResEmployee, error) {
	var employees []Employee
	var results []ResEmployee

	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	// Preload Payments and Advances filtered by last one month
	if err := db.
		Preload("Payments", "created_at >= ?", oneMonthAgo).
		Preload("SalaryAdvance", "created_at >= ?", oneMonthAgo).
		Find(&employees).Error; err != nil {
		return nil, err
	}

	// Build response
	for _, emp := range employees {
		var totalPayments float64
		var totalAdvances float64

		for _, pay := range emp.Payments {
			totalPayments += pay.Amount
		}

		for _, adv := range emp.SalaryAdvance {
			totalAdvances += adv.Amount
		}

		results = append(results, ResEmployee{
			Employees:     emp,
			TotalPayments: totalPayments,
			TotalAdvances: totalAdvances,
		})
	}

	return results, nil
}
