package controllers

import (
	"log"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/repositories"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)
func AddSupplyHandler(c *fiber.Ctx)error{
	
	newSupply, err := repositories.AddNewSupply(c)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to add supply",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Supply created successfully", newSupply)
}

//update supplies
func UpdateSupplyHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	supply := models.Supply{}
	if err := c.BodyParser(&supply); err != nil {
		return utils.NewErrorResponse(c,"failed to update",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	updatedSupply, err := repositories.UpdateSupply(c, id, &supply)
	log.Println(id)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to update",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Supply updated successfully", updatedSupply)
}
//delete supply
func DeleteSupplyHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	err := repositories.DeleteSupply(c, id)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to delete",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SendMessage(c,"Supply deleted successfully")
}

//get paginated supplies
func GetSuppliesHandler(c *fiber.Ctx)error{
	data, err := repositories.GetAllSupplies(c)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get supplies",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Supplies retrieved successfully", data)
}

func GetSupplyByIDHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	data, err := repositories.GetSupplyByID(c, id)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get supply",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Supply retrieved successfully", data)
}

func GetSupplierDebtsHandler(c *fiber.Ctx)error{
	data, err := repositories.GetSupplierDebts(db)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get supplier debts",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Supply debts retrieved successfully", data)
}