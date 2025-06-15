package controllers

import (
	"log"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewEmployeePayment(c *fiber.Ctx)error{
	p := models.Payment{}
	if err := c.BodyParser(&p); err != nil {
		log.Println(err.Error())
		return utils.NewErrorResponse(c,"failed to parse json data",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	
	employeePayment, err := models.AddEmployeePayment(&p)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to add employee payment",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Employee payment created successfully",employeePayment)
}

func UpdateEmployeePayment(c *fiber.Ctx)error{
	p := models.Payment{}
	if err := c.BodyParser(&p); err != nil {
		log.Println(err.Error())
		return utils.NewErrorResponse(c,"failed to parse json data",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	id, _ := uuid.Parse(c.Params("id"))
	
	employeePayment, err := models.UpdateEmployeePayment(&p, id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to update employee payment",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Employee payment updated successfully",employeePayment)
}

func DeleteEmployeePayment(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeletePayment(id)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to delete employee payment",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Employee payment deleted successfully",nil)
}

func GetRecentPaymentsHandler(c *fiber.Ctx)error{
	employeePayment, err := models.GetRecentPayments(10)
	if err != nil {
		return utils.NewErrorResponse(c,"failed to get employee payment",map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c,"Employee payment retrieved successfully",employeePayment)
}
