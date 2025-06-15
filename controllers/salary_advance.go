package controllers

import (
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateSalaryAdvanceHandler handles the creation of a new salary advance
func CreateSalaryAdvanceHandler(c *fiber.Ctx) error {
	SalaryAdvance := models.SalaryAdvance{}
	if err := c.BodyParser(&SalaryAdvance); err != nil {
		return utils.BadRequestResponse(c, "Failed to parse JSON data")
	}
	createdSalaryAdvance, err := models.CreateSalaryAdvance(c, &SalaryAdvance)
	if err != nil {
		return utils.BadRequestResponse(c, "Failed to create salary advance")
	}
	return utils.SuccessResponse(c, "Salary advance created successfully", createdSalaryAdvance)
}

// GetSalaryAdvanceByIDHandler handles the retrieval of a salary advance by ID
func GetSalaryAdvanceByIDHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	salaryAdvance, err := models.GetSalaryAdvanceByID(c, id)
	if err != nil {
		return utils.NotFoundResponse(c, "Salary advance not found")
	}
	return utils.SuccessResponse(c, "Salary advance retrieved successfully", salaryAdvance)
}

// UpdateSalaryAdvanceHandler handles the update of a salary advance
func UpdateSalaryAdvanceHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	var SalaryAdvance models.SalaryAdvance
	if err := c.BodyParser(&SalaryAdvance); err != nil {
		return utils.BadRequestResponse(c, "Failed to parse JSON data")
	}
	updatedSalaryAdvance, err := models.UpdateSalaryAdvance(c, id, &SalaryAdvance)
	if err != nil {
		return utils.BadRequestResponse(c, "Failed to update salary advance")
	}
	return utils.SuccessResponse(c, "Salary advance updated successfully", updatedSalaryAdvance)
}

// DeleteSalaryAdvanceHandler handles the deletion of a salary advance
func DeleteSalaryAdvanceHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeleteSalaryAdvance(c, id)
	if err != nil {
		return utils.BadRequestResponse(c, "Failed to delete salary advance")
	}
	return utils.SendMessage(c, "Salary advance deleted successfully")
}
