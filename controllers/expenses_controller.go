package controllers

import (
	"log"
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateExpensesHandler(c *fiber.Ctx) error {
	var expense models.Expenses
	if err := c.BodyParser(&expense); err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c, "Failed to parse JSON data", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	
	createdExpense, err := models.CreateExpenses(c, &expense)
	if err != nil {
		log.Println(err)
		return utils.NewErrorResponse(c, "Failed to create expense", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Expense created successfully", createdExpense)
}

//update expenses
func UpdateExpensesHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	expense := models.Expenses{}
	if err := c.BodyParser(&expense); err != nil {
		return utils.NewErrorResponse(c, "Failed to update", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	updatedExpense, err := models.UpdateExpenses(c, id, &expense)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to update", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Expense updated successfully", updatedExpense)
}

//delete expenses
func DeleteExpensesHandler(c *fiber.Ctx) error {
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeleteExpenses(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to delete", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Expense deleted successfully", nil)
}

//GET EXPENSES BY DATE
func GetExpensesByDateHandler(c *fiber.Ctx) error {
	expenses, err := models.GetExpensesByDate(c)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get expenses", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Expenses retrieved successfully", expenses)
}

//get expenses
func GetExpensesHandler(c *fiber.Ctx) error {
	expenses, err := models.GetExpenses(c)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get expenses", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Expenses retrieved successfully", expenses)
}

//GET EXPENSES OF A PROVIDED DURATION 
func GetExpensesByDurationHandler(c *fiber.Ctx) error {
	startDateStr := c.Params("start_date")
	endDateStr := c.Params("end_date")
	startDate, err := utils.ParseDate(startDateStr)
	if err != nil {
		return utils.NewErrorResponse(c, "Invalid start date format", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	endDate, err := utils.ParseDate(endDateStr)
	if err != nil {
		return utils.NewErrorResponse(c, "Invalid end date format", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	expenses, err := models.GetExpensesByDuration(c, startDate, endDate)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get expenses", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Expenses retrieved successfully", expenses)
}

func GetPaginatedExpensesByStation(c *fiber.Ctx)error{
	
	id, _ := uuid.Parse(c.Params("id"))
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	expenses, err := models.GetPaginatedExpensesByStation(c, id, page, pageSize)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get expenses", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Expenses retrieved successfully", expenses)
}