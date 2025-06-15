package controllers

import (
	"log"
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


func AddNewPumpHandler(c *fiber.Ctx) error {
	pump := new(models.Pump)
	if err := c.BodyParser(&pump); err != nil {
		log.Println(err.Error())
		return utils.NewErrorResponse(c, "Bad request", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	newPump, err := models.CreatePump(c, pump)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to add pump", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump created successfully", newPump)
}

//update Pump
func UpdatePumpHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	pump, err := models.UpdatePump(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to update pump", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump updated successfully", pump)
}

//delete pump
func DeletePumpHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeletePump(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to delete pump", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump deleted successfully", nil)
}

//asign pump to tank handler
func AssignPumpToTankHandler(c *fiber.Ctx) error {
	tankID, _ := uuid.Parse(c.Params("tank_id"))
	pumpID, _ := uuid.Parse(c.Params("pump_id"))
	err := models.AssignPumpToTank(c, tankID, pumpID)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to assign pump to tank", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump assigned to tank successfully", nil)
}

//Reassign pump to tank handler
func ReassignPumpToTankHandler(c *fiber.Ctx) error {
	tankID, _ := uuid.Parse(c.Params("tank_id"))
	pumpID, _ := uuid.Parse(c.Params("pump_id"))
	err := models.UnassignPumpFromTank(c, tankID, pumpID)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to reassign pump to tank", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump reassigned to tank successfully", nil)
}

//Get pump by id handler
func GetPumpByIDHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	pump, err := models.GetPumpByID(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get pump", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pump retrieved successfully", pump)
}

//Get pump by station handler
func GetPumpsByStationHandler(c *fiber.Ctx) error {
	stationID, _ := uuid.Parse(c.Params("id"))
	log.Println(stationID)
	pumps, err := models.GetPumpsByStation(c, stationID)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get pumps", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Pumps retrieved successfully", pumps)
}

