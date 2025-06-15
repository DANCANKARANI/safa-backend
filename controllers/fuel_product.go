package controllers

import (
	"log"
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateFuelProductHandler(c *fiber.Ctx) error {
	fuelProduct := models.FuelProduct{}
	if err := c.BodyParser(&fuelProduct); err != nil {
		log.Fatalln(err.Error())
		return utils.NewErrorResponse(c, "failed to create", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	createdFuelProduct, err := models.CreateFuelProduct(c, &fuelProduct)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to create", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Fuel product created successfully", createdFuelProduct)
}

func GetAllFuelProductsHandler(c *fiber.Ctx) error {
	fuelProducts, err := models.GetAllFuelProducts(c)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get fuel products",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Fuel products retrieved successfully", fuelProducts)
}

func GetFuelProductByIDHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	fuelProduct, err := models.GetFuelProductByID(c, id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to retrieve",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Fuel product retrieved successfully", fuelProduct)
}

func UpdateFuelProductHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	fuelProduct := models.FuelProduct{}
	if err := c.BodyParser(&fuelProduct); err != nil {
		return utils.NewErrorResponse(c,"failed to update", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	updatedFuelProduct, err := models.UpdateFuelProduct(c, id, &fuelProduct)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to update", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Fuel product updated successfully", updatedFuelProduct)
}

func DeleteFuelProductHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeleteFuelProduct(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "failed to delete", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SendMessage(c, "Fuel product deleted successfully")
}