package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (e *Expenses) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New()
	return nil
}

func (e *Expenses) BeforeUpdate(tx *gorm.DB) (err error) {
	e.UpdatedAt = time.Now()
	return nil
}

func CreateExpenses(c *fiber.Ctx, e *Expenses) (*Expenses, error) {
	db.AutoMigrate(&Expenses{})
	if err := db.Create(e).Error; err != nil {
		return nil, errors.New("failed to create expenses")
	}
	return e, nil
}

func UpdateExpenses(c *fiber.Ctx, id uuid.UUID, updatedData *Expenses) (*Expenses, error) {
	var expenses Expenses
	if err := db.First(&expenses, "id = ?", id).Error; err != nil {
		return nil, errors.New("expenses not found")
	}
	expenses.ExpenseType = updatedData.ExpenseType
	if updatedData.Amount != 0 {
		expenses.Amount = updatedData.Amount
	}
	if updatedData.Description != "" {
		expenses.Description = updatedData.Description
	}
	expenses.ExpenseDate = updatedData.ExpenseDate

	if err := db.Save(&expenses).Error; err != nil {
		return nil, errors.New("failed to update expenses")
	}
	return &expenses, nil
}

//get expenses
func GetExpenses(c *fiber.Ctx) ([]Expenses, error) {
	var expenses []Expenses
	if err := db.Find(&expenses).Error; err != nil {
		return nil, errors.New("failed to get expenses")
	}
	return expenses, nil
}

//get expenses by id
func GetExpensesByID(c *fiber.Ctx, id uuid.UUID) (*Expenses, error) {
	var expenses Expenses
	if err := db.First(&expenses, "id = ?", id).Error; err != nil {
		return nil, errors.New("expenses not found")
	}
	return &expenses, nil
}
//get expenses by date
type ResExpenses struct {
	Expenses []Expenses
	Total    float64
}

func GetExpensesByDate(c *fiber.Ctx) (*ResExpenses, error) {
	// Get date query or default to current day
	dateParam := c.Query("date")
	var targetDate time.Time
	var err error

	if dateParam == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD")
		}
	}

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var expenses []Expenses
	if err := db.Where("expense_date >= ? AND expense_date < ?", startOfDay, endOfDay).Find(&expenses).Error; err != nil {
		return nil, errors.New("failed to get expenses")
	}

	var total float64
	for _, exp := range expenses {
		total += exp.Amount
	}

	return &ResExpenses{
		Expenses: expenses,
		Total:    total,
	}, nil
}

//delete expenses
func DeleteExpenses(c *fiber.Ctx, id uuid.UUID) error {
	db.AutoMigrate(&Expenses{})
	var expenses Expenses
	if err := db.First(&expenses, "id = ?", id).Error; err != nil {
		return errors.New("expenses not found")
	}
	if err := db.Delete(&expenses).Error; err != nil {
		return errors.New("failed to delete expenses")
	}
	return nil
}

//GET EXPENSES OF A PROVIDED DURATION 
func GetExpensesByDuration(c *fiber.Ctx, startDate, endDate time.Time) ([]Expenses, error) {
	var expenses []Expenses

	// If both dates are zero (not provided), use current month's range
	if startDate.IsZero() && endDate.IsZero() {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	}

	if err := db.Where("expense_date BETWEEN ? AND ?", startDate, endDate).Find(&expenses).Error; err != nil {
		return nil, errors.New("expenses not found")
	}

	return expenses, nil
}

//get expenses by duration

func GetPaginatedExpensesByStation(c *fiber.Ctx, stationID uuid.UUID, page, pageSize int) ([]Expenses, error) {
	var expenses []Expenses
	if err := db.Where("station_id = ?", stationID).
		Order("expense_date DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&expenses).Error; err != nil {
		return nil, errors.New("failed to get expenses")
	}
	return expenses, nil
}