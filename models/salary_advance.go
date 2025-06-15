package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateSalaryAdvance(c *fiber.Ctx, e *SalaryAdvance) (*SalaryAdvance, error) {
	e.ID = uuid.New()
	db.AutoMigrate(&SalaryAdvance{})
	if err := db.Create(e).Error; err != nil {
		return nil, errors.New("failed to create salary advance")
	}
	return e, nil
}

func UpdateSalaryAdvance(c *fiber.Ctx, id uuid.UUID, updatedData *SalaryAdvance) (*SalaryAdvance, error) {
	var salaryAdvance SalaryAdvance
	if err := db.First(&salaryAdvance, "id = ?", id).Error; err != nil {
		return nil, errors.New("salary advance not found")
	}
	salaryAdvance.EmployeeID = updatedData.EmployeeID
	salaryAdvance.Amount = updatedData.Amount
	if err := db.Save(&salaryAdvance).Error; err != nil {
		return nil, errors.New("failed to update salary advance")
	}
	return &salaryAdvance, nil
}

func DeleteSalaryAdvance(c *fiber.Ctx, id uuid.UUID) error {
	var salaryAdvance SalaryAdvance
	if err := db.First(&salaryAdvance, "id = ?", id).Error; err != nil {
		return errors.New("salary advance not found")
	}
	if err := db.Where("id = ?", id).Delete(&salaryAdvance).Error; err != nil {
		return errors.New("failed to delete salary advance")
	}
	return nil
}

func GetSalaryAdvanceByID(c *fiber.Ctx, id uuid.UUID) (*SalaryAdvance, error) {
	var salaryAdvance SalaryAdvance
	if err := db.First(&salaryAdvance, "id = ?", id).Error; err != nil {
		return nil, errors.New("salary advance not found")
	}
	return &salaryAdvance, nil
}