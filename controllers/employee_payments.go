package controllers

import (
	"log"

	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewEmployeePayment(c *fiber.Ctx) error {
	var p models.Payment

	// Parse JSON body
	if err := c.BodyParser(&p); err != nil {
		log.Println("JSON parsing error:", err)
		return utils.NewErrorResponse(c, "Invalid JSON payload", map[string][]string{
			"error": {err.Error()},
		}, fiber.StatusBadRequest)
	}
	p.PaidMonth = c.Query("month") 
	// Validate required fields manually
	if p.EmployeeID == uuid.Nil || p.PaidMonth == "" || p.Amount <= 0 {
		return utils.NewErrorResponse(c, "Missing required payment fields", map[string][]string{
			"employee_id": {"Employee ID is required"},
			"paid_month":  {"Paid month is required in YYYY-MM format"},
			"amount":      {"Amount must be greater than 0"},
		}, fiber.StatusBadRequest)
	}

	// Save the payment
	employeePayment, err := models.AddEmployeePayment(c,&p)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to add employee payment", map[string][]string{
			"error": {err.Error()},
		}, fiber.StatusBadRequest)
	}

	// Return success response
	return utils.SuccessResponse(c, "Employee payment created successfully", employeePayment)
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

//get payroll
func GetPayrollReportHandler(c *fiber.Ctx) error {
	month := c.Query("month")

	report, err := models.GetReport(month)
	if err != nil {
		utils.NewErrorResponse( c, "Failed to get report", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}

	return utils.SuccessResponse(c, "Report retrieved successfully", report)
}

