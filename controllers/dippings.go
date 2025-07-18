package controllers

import (
	"log"
	"time"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateDippingHandler(c *fiber.Ctx)error {
	var dipping models.Dippings
	if err := c.BodyParser(&dipping); err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c,"failed to parse json data",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	dipping.ID = uuid.New()
	createdDipping, err := models.CreateDipping(&dipping)
	if err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c,"failed to create",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dipping created successfully",createdDipping)
}

//get all dippings
func GetAllDippingsHandler(c *fiber.Ctx)error {
	data, err := models.GetAllDippings(c)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get dippings",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dippings retrieved successfully",data)
}

//get dipping by id
func GetDippingByIDHandler(c *fiber.Ctx)error {
	id, _ := uuid.Parse(c.Params("id"))
	data, err := models.GetDippingByID(id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get dipping",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dipping retrieved successfully",data)
}

//get dippings by station
func GetDippingsByStationHandler(c *fiber.Ctx)error {
	id, _ := uuid.Parse(c.Params("station_id"))
	data, err := models.GetDippingByStationID(id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get dippings",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dippings retrieved successfully",data)
}



func GetDippingByDippingDateHandler(c *fiber.Ctx)error {
	dateStr := c.Params("date")
	date, err := time.Parse("2006-01-02", dateStr) // adjust format as needed
	if err != nil {
		return utils.NewErrorResponse(c,"invalid date format",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	data, err := models.GetDippingByDippingDate(date)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get dippings",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dippings retrieved successfully",data)
}

//get dipping by fuel product
func GetDippingByFuelProductHandler(c *fiber.Ctx)error {
	id, _ := uuid.Parse(c.Params("id"))
	data, err := models.GetDippingByFuelProductID(id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get dippings",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dippings retrieved successfully",data)
}

//update dipping
func UpdateDippingHandler(c *fiber.Ctx)error {
	id, _ := uuid.Parse(c.Params("id"))
	
	updatedDipping, err := models.UpdateDippings(c, id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to update",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dipping updated successfully",updatedDipping)
}

//compare dippings and sales
func CompareDippingsAndSales(c *fiber.Ctx)error {
	data, err := models.ComparePumpReadingsWithDippings(c)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get dippings",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Dippings retrieved successfully",data)
}

//get openingDippings handler
func GetOpeningDipHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("tank_id"))
	openingDipping, err := models.GetOpeningDippings(c, id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get opening dippings",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Opening dippings retrieved successfully",openingDipping)
}

//get corresponding closing sales for the dippings
func GetClosingSalesHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("tank_id"))
	closingDipping, err := models.GetLatestReadingsSumByTankID(id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get closing sales",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Closing corresponding sales retrieved successfully",closingDipping)
}