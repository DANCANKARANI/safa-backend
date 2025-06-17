package controllers

import (
	"github.com/dancankarani/safa/repositories"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
)

func GetFuelStockHandler(c *fiber.Ctx) error {
	fuelStock,err := repositories.GetFuelStocksWithDetails()
	if err != nil{
		return utils.NewErrorResponse(c, "failed to get fuel stock", map[string][]string{"errors": {err.Error()}}, fiber.StatusInternalServerError)
	}
	return utils.SuccessResponse(c, "fuel stock retrieved successfully", fuelStock)
}