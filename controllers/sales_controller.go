package controllers

import (
	"log"
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewSalesHandler(c *fiber.Ctx) error {
	createdSales, err := models.AddNewSales(c)
	if err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c, "Failed to add sales", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Sales added successfully", createdSales)
}

func UpdateSalesHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	sales := models.Sales{}
	if err := c.BodyParser(&sales); err != nil {
		return utils.NewErrorResponse(c, "Failed to update", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	updatedSales, err := models.UpdateSales(c, id, &sales)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to update", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Sales updated successfully", updatedSales)
}

func GetSalesHandler(c *fiber.Ctx) error {
	sales, err := models.GetSales(c)
	if err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c, "Failed to get sales", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Sales retrieved successfully", sales)
}

func GetSalesByDateHandler(c *fiber.Ctx) error {
	start := c.Params("date")
	end := c.Params("date")
	start_date, err := utils.ParseDate(start)
	if err != nil {
		return utils.NewErrorResponse(c, "Invalid date format", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	end_date, err := utils.ParseDate(end)
	if err != nil {
		return utils.NewErrorResponse(c, "Invalid date format", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	
	sales, err := models.GetSalesByDate(c, start_date, end_date)
	if err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c, "Failed to get sales", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Sales retrieved successfully", sales)
}

func DeleteSalesHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeleteSales(c, id)
	if err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c, "Failed to delete sales", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Sales deleted successfully", nil)
}