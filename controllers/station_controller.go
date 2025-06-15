package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
)

func NewStationHandler(c *fiber.Ctx) error{
	station := models.Station{}
	if err := c.BodyParser(&station); err != nil {
		return utils.NewErrorResponse(c,"failed to parse json data",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}

	newStation, err := models.AddNewStation(c, station)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to add station",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Station created successfully", newStation)
}

func ReadAllStationsController(c *fiber.Ctx) error{
	stations, err := models.GetAllStations(c)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get all stations",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Stations retrieved successfully",stations)
}

func ReadStationByIDController(c *fiber.Ctx) error{
	id, _ := uuid.Parse(c.Params("id"))
	station, err := models.GetStationByID(c, id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get station by id",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Station retrieved successfully",station)
}

func GetAllSalesInStationHandler(c *fiber.Ctx)error{
	
	sales, err := models.GetSummationOfSalesAndLiters(c)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get sales",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Sales retrieved successfully", sales)
}

func GetStationExpensesHandler(c *fiber.Ctx)error{
	
	expenses, err := models.GetSummationOfExpenses(c)
	if err != nil{
		return utils.NewErrorResponse(c,"failed to get expenses",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Expenses retrieved successfully", expenses)
}