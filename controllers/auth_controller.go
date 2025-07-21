package controllers

import (
	"log"
	"time"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/repositories"
	"github.com/dancankarani/safa/services"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
)

type ResponseLogin struct {
	Employee models.Employee `json:"employee"`
	Token    string           `json:"token"`
}

// Login handles the login request for employees
func LoginUser(c *fiber.Ctx) error {
	employee := models.Employee{}
	err_str := map[string][]string{}
	if err := c.BodyParser(&employee); err != nil {
		log.Println("Error parsing request body:", err.Error())
		err_str["error"] = []string{"Invalid request body"}
		return utils.NewErrorResponse(c, "Invalid request body", err_str, fiber.StatusBadRequest)
	}

	if employee.Email == "" || employee.Password == "" {
		err_str["error"] = []string{"Email and password are required"}
		return utils.NewErrorResponse(c, "Email and password are required", err_str, fiber.StatusBadRequest)
	}

	//get the employee by email
	existingEmployee, err := repositories.GetEmployeeByEmail(employee.Email)
	if err != nil {
		err_str["error"] = []string{"invalid credentials"}
		return utils.NewErrorResponse(c, "invalid credentials", err_str, fiber.StatusNotFound)
	}

	err = services.CompareHashAndPassword(existingEmployee.Password, employee.Password)
	if err != nil {
		err_str["error"] = []string{"invalid credentials"}
		return utils.NewErrorResponse(c, "invalid credentials", err_str, fiber.StatusUnauthorized)
	}

	// Generate JWT token
	exp := time.Hour * 24
	full_name := existingEmployee.FirstName + " " + existingEmployee.LastName
	token, err := services.GenerateToken(services.Claims{
		UserID:   &existingEmployee.ID,
		FullName: full_name,
		Role:     existingEmployee.Role,
	}, exp)
	if err != nil {
		err_str["error"] = []string{"failed to generate token"}
		return utils.NewErrorResponse(c, "failed to generate token", err_str, fiber.StatusInternalServerError)
	}

	//set the token
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true, 
		Secure:   true,
    	SameSite: "None",
		Path:     "/",
	})
	return utils.SuccessResponse(c, "Login successful", ResponseLogin{
		Employee: existingEmployee,
		Token:    token,
	})
}
