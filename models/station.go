package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddNewStation (c *fiber.Ctx, station Station) (*Station,error) {
	station.ID = uuid.New()
	db.AutoMigrate(&Station{})
	if err := db.Create(&station).Error; err != nil {
		return nil, errors.New("failed to add new station")
	}
	return &station, nil	
}

func GetStationByID(c* fiber.Ctx, id uuid.UUID) (*Station, error) {
	var station Station
	if err := db.Preload("Expenses").Preload("Tanks.Pumps.Sales").Where("id = ?", id).First(&station).Error; err != nil {
		return nil, errors.New("failed to get station by id")
	}
	return &station, nil
}

func GetAllStations(c *fiber.Ctx) (*[]Station, error) {
	var stations []Station
	if err := db.
		Preload("Tanks").
		Preload("Expenses").
		Preload("Tanks.Pumps.Nozzles").
		Preload("Tanks.Dippings").
		Preload("Tanks.Pumps.Sales").
		Find(&stations).Error; err != nil {
		return nil, errors.New("failed to get all stations")
	}
	return &stations, nil
}

// GetSummationOfSalesAndLiters returns the total sales amount and dispersed liters from all pump readings in every station
type ResStationSales struct {
	StationID   uuid.UUID `json:"station_id"`
	StationName string    `json:"station_name"`
	TotalSales  float64   `json:"total_sales"`
	TotalLiters float64   `json:"total_liters"`
	MpesaAmount float64   `json:"mpesa_amount"`
	BankDeposit float64   `json:"bank_deposit"`
}

func GetSummationOfSalesAndLiters(c *fiber.Ctx) (*[]ResStationSales, error) {
	// Parse optional `date` query param (format: YYYY-MM-DD)
	dateParam := c.Query("date")
	var targetDate time.Time
	var err error

	if dateParam == "" {
		// Use current date if no date is provided
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD")
		}
	}

	startOfDay := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var stations []Station
	if err := db.
		Preload("Tanks.Pumps.Readings").
		Find(&stations).Error; err != nil {
		return nil, errors.New("failed to get all stations")
	}

	var results []ResStationSales
	for _, station := range stations {
		var totalSales float64
		var totalLiters float64
		var mpesaAmount float64
		var bankDeposit float64

		for _, tank := range station.Tanks {
			for _, pump := range tank.Pumps {
				for _, reading := range pump.Readings {
					// Only include readings within the target day
					if reading.ReadingDate.After(startOfDay) && reading.ReadingDate.Before(endOfDay) {
						totalSales += reading.TotalSalesAmount
						totalLiters += reading.LitersDispensed
						mpesaAmount += reading.MpesaAmount
						bankDeposit += reading.BankDeposit
					}
				}
			}
		}

		results = append(results, ResStationSales{
			StationID:   station.ID,
			StationName: station.Name,
			TotalSales:  totalSales,
			TotalLiters: totalLiters,
			MpesaAmount: mpesaAmount,
			BankDeposit: bankDeposit,
		})
	}

	return &results, nil
}


//get station expenses and totals

type StationExpenses struct {
	StationID   uuid.UUID `json:"station_id"`
	StationName string    `json:"station_name"`
	Total  float64   `json:"total_expenses"`
}

type ResSationExpenses struct{
	StationExpenses []StationExpenses `json:"station_expenses"`
	TotalExpenses float64 `json:"total_expenses"`
}

func GetSummationOfExpenses(c *fiber.Ctx) (*ResSationExpenses, error) {
	dateParam := c.Query("date")
	var targetDate time.Time
	var err error

	if dateParam == "" {
		now := time.Now().UTC() // or .Local() if DB stores local time
		// Set targetDate to today at 8:00 AM
		targetDate = time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
	} else {
		targetDate, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD")
		}
		// Adjust targetDate to 8:00 AM of that date
		targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 8, 0, 0, 0, targetDate.Location())
	}

	startOfDay := targetDate
	endOfDay := startOfDay.Add(24 * time.Hour) // 8 AM next day

	var stations []Station
	if err := db.Preload("Expenses").Find(&stations).Error; err != nil {
		return nil, errors.New("failed to get all stations")
	}

	var results []StationExpenses
	var totalExpenses float64

	for _, station := range stations {
		var total float64
		for _, expense := range station.Expenses {
			// Include expenses with ExpenseDate >= startOfDay AND < endOfDay
			if !expense.ExpenseDate.Before(startOfDay) && expense.ExpenseDate.Before(endOfDay) {
				total += expense.Amount
			}
		}

		results = append(results, StationExpenses{
			StationID:   station.ID,
			StationName: station.Name,
			Total:       total,
		})
		totalExpenses += total
	}

	return &ResSationExpenses{
		StationExpenses: results,
		TotalExpenses:   totalExpenses,
	}, nil
}

