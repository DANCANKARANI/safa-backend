package controllers

import (
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/repositories"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewPumpReadings(c *fiber.Ctx) error {
	pumpReadings := models.PumpReadings{}
	if err := c.BodyParser(&pumpReadings); err != nil {
		return utils.NewErrorResponse(c, "Failed to parse JSON data", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	newPumpReadings, err := repositories.AddPumpReadings(c, pumpReadings)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to add pump readings", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump readings added successfully", newPumpReadings)
}




func UpdatePumpReadingsHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	pumpReadings := models.PumpReadings{}
	if err := c.BodyParser(&pumpReadings); err != nil {
		return utils.NewErrorResponse(c, "Failed to parse JSON data", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	updatedPumpReadings, err := models.UpdatePumpReadings(c, id, pumpReadings)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to update pump readings", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump readings updated successfully", updatedPumpReadings)
}

func DeletePumpReadingsHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeletePumpReadings(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to delete pump readings", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump readings deleted successfully", nil)
}

func GetPumpReadingsHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	pumpReadings, err := models.GetLatestPumpReadingsByStationID(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get pump readings", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump readings retrieved successfully", pumpReadings)
}

type PumpReadingsResponse struct {
	PumpReadings []models.PumpReadings `json:"pumpReadings"`
	Total        int64                 `json:"total"`
}

func GetOrderedPumpReadingsHandler(c *fiber.Ctx) error {
	pumpReadings, total, err := models.GetPaginatedPumpReadings(c, 1, 10)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get pump readings", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	pumpReadingsResponse := PumpReadingsResponse{PumpReadings: pumpReadings, Total: total}

	return utils.SuccessResponse(c, "Pump readings retrieved successfully", pumpReadingsResponse)
}


//get all sales by date

func GetAllSalesByDateHandler(c *fiber.Ctx) error {
	
	sales, err := models.GetTotalSalesByDate(c)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get sales", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}

	return utils.SuccessResponse(c, "Sales retrieved successfully", sales)
}
