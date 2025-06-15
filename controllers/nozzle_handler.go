package controllers 

import (
	"log"
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNozzleHandler(c *fiber.Ctx) error{
	nozzle := new(models.Nozzle)
	if err := c.BodyParser(nozzle); err != nil {
		log.Println(err.Error())
		return utils.NewErrorResponse(c, "Bad request", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	
	newNozzle, err := models.CreateNozzle(c, *nozzle)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to add nozzle", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Nozzle created successfully", newNozzle)
}

func UpdateNozzleHandler(c *fiber.Ctx)error{
	nozzle := new(models.Nozzle)
	if err := c.BodyParser(nozzle); err != nil {
		log.Println(err.Error())
		utils.NewErrorResponse(c, "Bad request", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	id, _:= uuid.Parse(c.Params("id"))
	resp, err := models.UpdateNozzle(c, *nozzle,id)
	if err != nil {
		log.Println(err.Error())
		err_msg := map[string][]string{
			"errors":{err.Error()},
		}
		return utils.NewErrorResponse(c,"bad request",err_msg, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Nozzle updated successfully", resp)
}

func DeleteNozzleHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeleteNozzle(id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to delete nozzle", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Nozzle deleted successfully", nil)
}

func GetAllNozzlesHandler(c *fiber.Ctx)error{
	nozzles, err := models.GetAllNozzles()
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get nozzles", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Nozzles retrieved successfully", nozzles)
}

func GetNozzleByIDHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	nozzle, err := models.GetNozzleByID(id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get nozzle", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Nozzle retrieved successfully", nozzle)
}

