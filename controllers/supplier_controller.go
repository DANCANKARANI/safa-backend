package controllers

import (
	"github.com/dancankarani/safa/database"
	"github.com/dancankarani/safa/models"
	"github.com/dancankarani/safa/repositories"
	"github.com/dancankarani/safa/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddSupplierHandler(c *fiber.Ctx) error {
	s := &models.Supplier{}
	if err := c.BodyParser(s); err != nil {
		return utils.NewErrorResponse(c, "Failed to parse JSON data", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	s, err := models.AddSupplier(c, s)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to add supplier", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Supplier added successfully", s)
}

func UpdateSupplierHandler(c *fiber.Ctx) error {
	id,_ := uuid.Parse(c.Params("id"))
	s := models.Supplier{}
	if err := c.BodyParser(&s); err != nil {
		return utils.NewErrorResponse(c, "Failed to parse JSON data", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	updatedData, err := models.UpdateSupplier(c,id, &s)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to update supplier", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Supplier updated successfully", updatedData)
}

func GetSuppliersHandler(c *fiber.Ctx) error {

	data, err := models.GetAllSuppliers(c)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get suppliers", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Suppliers retrieved successfully", data)
}

func GetSupplierHandler(c *fiber.Ctx) error {
	id,_ := uuid.Parse(c.Params("id"))
	data, err := models.GetSupplierByID(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get supplier", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Supplier retrieved successfully", data)
}

func DeleteSupplierHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	err := models.DeleteSupplier(c, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to delete supplier", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SendMessage(c, "Supplier deleted successfully")
}

var db = database.ConnectDB()
type GetSupplierBalance struct{
	SupplierID uuid.UUID `json:"supplier_id"`
	DebtYouOwe float64 `json:"debt_you_owe"`
	CreditTheyOwe float64 `json:"credit_they_owe"`
	NetBalance float64 `json:"net_balance"`
}
func GetSupplierBalanceHandler(c *fiber.Ctx)error{
	id, _ := uuid.Parse(c.Params("id"))
	debt_you_owe, credit_they_owe, net_balance, err := repositories.GetSupplierBalance(db, id)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get supplier balance", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	data := GetSupplierBalance{SupplierID: id, DebtYouOwe: debt_you_owe, CreditTheyOwe: credit_they_owe, NetBalance: net_balance}
	return utils.SuccessResponse(c, "Supplier balance retrieved successfully", data)
}

//get all balances handler
func GetAllSupplierBalancesHandler(c *fiber.Ctx)error{
	data, err := repositories.GetSupplierDebts( db)
	if err != nil {
		return utils.NewErrorResponse(c, "Failed to get supplier balances", map[string][]string{"error": {err.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Supplier balances retrieved successfully", data)
}

//add supplier payments
func AddSupplierPayments(c *fiber.Ctx)error{
	payment, error := repositories.AddSupplierPayments(c)
	if error != nil {
		return utils.NewErrorResponse(c, "Failed to add supplier payment", map[string][]string{"error": {error.Error()}}, fiber.StatusBadRequest)
	}
	return utils.SuccessResponse(c, "Supplier payment added successfully", payment)
}