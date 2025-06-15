package controllers

import (
	"time"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewPumpReadings(c *fiber.Ctx) error {
	pumpReadings := models.PumpReadings{}
	if err := c.BodyParser(&pumpReadings); err != nil {
		return utils.NewErrorResponse(c, "Failed to parse JSON data", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	newPumpReadings, err := models.AddPumpReadings(c, pumpReadings)
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
	start_date := c.Params("start_date")
	start, err := utils.ParseDate(start_date)
	if err != nil {
		return utils.NewErrorResponse(c, "Invalid date format", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}

	end_date := c.Params("end_date")
	end, err := utils.ParseDate(end_date)
	if err != nil {
		return utils.NewErrorResponse(c, "Invalid date format", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}

	// Normalize start to 00:00:00 and end to 23:59:59.999
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), end.Location())

	sales, err := models.GetTotalSalesByDate(c, start, end)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get sales", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}

	return utils.SuccessResponse(c, "Sales retrieved successfully", sales)
}
