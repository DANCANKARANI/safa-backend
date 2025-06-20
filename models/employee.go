package models

import (
	"errors"
	"fmt"
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
		log.Println(err.Error())
		return nil, errors.New("failed to create employee")
	}

	//send Login details to the user
	htmlBody := fmt.Sprintf(`
	<html>
	<body style="background-color: #f0f0f0; margin: 0; padding: 0;">
		<div style="max-width: 600px; margin: 40px auto; background-color: #ffffff; padding: 30px; border-radius: 10px; font-family: Arial, sans-serif; color: #333;">
		
		<h1 style="font-size: 26px; font-weight: bold; color: #2c3e50; margin-top: 0; margin-bottom: 20px;">
			SAFA Login Credentials
		</h1>
		
		<p style="font-size: 16px; margin-bottom: 20px;">
			Dear <strong>%s</strong>,
		</p>
		
		<p style="font-size: 16px; margin-bottom: 10px;">
			Your account has been created successfully. Below are your login details:
		</p>
		
		<p style="font-size: 16px; margin-bottom: 10px;">
			<strong>Username:</strong> <span style="background-color: #eef; padding: 4px 8px; border-radius: 4px;">%s</span>
			<strong>Password:</strong> <span style="background-color: #eef; padding: 4px 8px; border-radius: 4px;">%s</span>
		</p>
		
		<p style="font-size: 14px; color: #777; margin-top: 30px;">
			Please change your password after your first login for security purposes.
		</p>
		
		<p style="font-size: 14px; color: #aaa; margin-top: 40px;">
			&copy; 2025 SAFA Systems
		</p>
		
		</div>
	</body>
	</html>
	`, e.FirstName,e.Email, password)

	go services.SendEmail(e.Email,"Password", htmlBody)
	log.Println("Pasword:", password)
	return e, nil
}

func GetEmployeeByID(c *fiber.Ctx, id uuid.UUID) (*Employee, error) {
	var employee Employee
	if err := db.Preload("Station").
		Preload("Payments").
		Preload("SalaryAdvance").

		First(&employee, "id = ?", id).Error; err != nil {
		return nil, errors.New("Employee not found")
	}
	return &employee, nil
}
type ResEmployees struct {
	Employees []Employee `json:"employees"`
	Count     int        `json:"count"`
}

func GetAllEmployees(c *fiber.Ctx) (*ResEmployees, error) {
	var employees []Employee
	if err := db.Preload("Station").
		Preload("Payments").
		Preload("SalaryAdvance").
		Find(&employees).Error; err != nil {
		return nil, errors.New("failed to retrieve employees")
	}

	return &ResEmployees{
		Employees: employees,
		Count:     len(employees),
	}, nil
}


func UpdateEmployee(c *fiber.Ctx, id uuid.UUID, updatedData *Employee) (*Employee, error) {
	var employee Employee
	if err := db.First(&employee, "id = ?", id).Error; err != nil {
		return nil, errors.New("employee not found")
	}

	// Only update if field is non-empty or meaningful
	if updatedData.FirstName != "" {
		employee.FirstName = updatedData.FirstName
	}

	if updatedData.LastName != "" {
		employee.LastName = updatedData.LastName
	}

	if updatedData.Position != "" {
		employee.Position = updatedData.Position
	}

	if updatedData.PhoneNumber != "" {
		if !services.ValidatePhoneNumber(updatedData.PhoneNumber) {
			return nil, errors.New("invalid phone number")
		}
		employee.PhoneNumber = updatedData.PhoneNumber
	}

	if updatedData.Email != "" {
		if !services.ValidateEmail(updatedData.Email) {
			return nil, errors.New("invalid email")
		}
		employee.Email = updatedData.Email
	}

	if updatedData.Password != "" {
		hashedPassword, err := services.HashPassword(updatedData.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		employee.Password = hashedPassword
	}

	if updatedData.Salary != 0 {
		employee.Salary = updatedData.Salary
	}

	// Save updates
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
