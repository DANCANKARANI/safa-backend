package controllers

import (
	"log"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewTankHandler(c *fiber.Ctx)error{
	tank := new(models.Tank)
	
	if err := c.BodyParser(tank); err != nil {
		log.Println(err.Error())
		err_msg := map[string][]string{
			"errors":{err.Error()},
		}
		return utils.NewErrorResponse(c,"bad request",err_msg, fiber.StatusBadRequest)
	}
	resp, err := models.CreateTank(c, *tank)
	if err != nil {
		log.Println(err.Error())
		err_msg := map[string][]string{
			"errors":{err.Error()},
		}
		return utils.NewErrorResponse(c,"bad request",err_msg, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Tank created successfully", resp)
}
//get all tanks in a station
func GetAllTanksHandler(c *fiber.Ctx)error{
	station_id,err := uuid.Parse(c.Params("id"))
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get tanks",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	tanks, err := models.GetTanksByStation(c,station_id)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get tanks",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Tanks retrieved successfully", tanks)
}

//Get tank by id
func GetTankByIDHandler(c *fiber.Ctx)error{
	id, _:= uuid.Parse(c.Params("id"))
	tank, err := models.GetTankByID(c, id)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get tank",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Tank retrieved successfully", tank)
}

//update Tank handler
func UpdateTankHandler(c *fiber.Ctx)error{
	tank := new(models.Tank)
	
	if err := c.BodyParser(tank); err != nil {
		log.Println(err.Error())
		err_msg := map[string][]string{
			"errors":{err.Error()},
		}
		return utils.NewErrorResponse(c,"bad request",err_msg, fiber.StatusBadRequest)
	}
	id, _:= uuid.Parse(c.Params("id"))
	resp, err := models.UpdateTank(c, *tank,id)
	if err != nil {
		log.Println(err.Error())
		err_msg := map[string][]string{
			"errors":{err.Error()},
		}
		return utils.NewErrorResponse(c,"bad request",err_msg, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Tank updated successfully", resp)
}

//delete Tank handler
func DeleteTankHandler(c *fiber.Ctx)error{
	id, _:= uuid.Parse(c.Params("id"))
	err := models.DeleteTank(c, id)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to delete",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SendMessage(c,"Tank deleted successfully")
}
