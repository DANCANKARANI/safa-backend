package controllers

import (
	"log"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/repositories"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateEmployee handles the creation of a new employee
func CreateEmployee(c *fiber.Ctx) error {
	var employee models.Employee
	err_str := map[string][]string{}
	if err := c.BodyParser(&employee); err != nil {
		log.Println("failed to parse json data", err.Error())
		err_str["error"] = []string{"failed to parse json data"}
		return utils.NewErrorResponse(c, "failed to parse json data", err_str, fiber.StatusBadRequest)
	}
	//check if employee with this email alredy exist
	exists, err := repositories.EmployeeExists(employee.Email)
	if err != nil {
		err_str["error"] = []string{err.Error()}
		return utils.NewErrorResponse(c, "failed to check employee existence", err_str, fiber.StatusBadRequest)
	}
	if exists {
		err_str["error"] = []string{"employee with this email already exist"}
		return utils.NewErrorResponse(c, "employee with this email already exist", err_str, fiber.StatusBadRequest)
	}
	if _, err := models.CreateEmployee(c, &employee); err != nil {
		log.Println("failed to add employee", err.Error())
		err_str["error"] = []string{err.Error()}
		return utils.NewErrorResponse(c, "failed to add employee", err_str, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Employee created successfully", employee)
	
}

// GetEmployeeByID retrieves an employee by their ID
func GetEmployeeByID(c *fiber.Ctx) error {
	id,_ := uuid.Parse(c.Params("id"))
	employee, err := models.GetEmployeeByID(c, id)
	if err != nil {
		return utils.NotFoundResponse(c, "Employee not found")
	}
	return utils.SuccessResponse(c, "Employee retrieved successfully", employee)
}

// GetAllEmployees retrieves all employees
func GetAllEmployees(c *fiber.Ctx) error {
	employees, err := models.GetAllEmployees(c)
	err_str := map[string][]string{}
	if err != nil {
		err_str["error"] = []string{err.Error()}
		return utils.NewErrorResponse(c,"failed to retrieve employees",err_str, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Employees retrieved successfully", employees)
}

// UpdateEmployee updates an existing employee's details
func UpdateEmployee(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	var employee models.Employee
	err_str := map[string][]string{}
	if err := c.BodyParser(&employee); err != nil {
		err_str["error"] = []string{"failed to parse json data"}
		return utils.NewErrorResponse(c, "failed to update employee", err_str, fiber.StatusBadRequest)
	}
	employee.ID = id // Ensure the ID is set for the update
	updated, err := models.UpdateEmployee(c, id, &employee);
	if err != nil {
		err_str["error"] = []string{err.Error()}
		return utils.NewErrorResponse(c, "failed to update employee", err_str, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Employee updated successfully", updated)
}

// DeleteEmployee deletes an employee by their ID
func DeleteEmployee(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	err_str := map[string][]string{}
	if err := models.DeleteEmployee(c, id); err != nil {
		err_str["error"] = []string{err.Error()}
		return utils.NewErrorResponse(c, "failed to delete employee", err_str, fiber.StatusBadRequest)
	}
	return utils.SendMessage(c, "Employee deleted successfully")
}

func GetEmployeePaymentsAndAdvancesHandler(c *fiber.Ctx)error{
	employeePayment, err := models.GetEmployeePaymentsAndAdvances()
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get employee payment",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Employee payment retrieved successfully",employeePayment)
}