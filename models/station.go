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
		Preload("StationFuelProducts.FuelProduct"). // 👈 Preload fuel prices and names
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
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			return nil, fmt.Errorf("invalid date format. Use YYYY-MM-DD")
		}
	}

	// Normalize targetDate to midnight for equality checks
	targetDate = time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())

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
					if reading.BusinessDay.Equal(targetDate) {
						totalSales += reading.TotalSalesAmount
						totalLiters += reading.LitersDispensed
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
func GetMonthlySummationOfExpenses(c *fiber.Ctx) (*ResSationExpenses, error) {
	monthParam := c.Query("month")
	var startOfMonth time.Time
	var err error

	if monthParam == "" {
		now := time.Now().UTC() // or .Local() if DB uses local time
		startOfMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	} else {
		startOfMonth, err = time.Parse("2006-01", monthParam)
		if err != nil {
			return nil, fmt.Errorf("invalid month format. Use YYYY-MM")
		}
	}

	// First day of next month
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var stations []Station
	if err := db.Preload("Expenses").Find(&stations).Error; err != nil {
		return nil, errors.New("failed to get all stations")
	}

	var results []StationExpenses
	var totalExpenses float64

	for _, station := range stations {
		var total float64
		for _, expense := range station.Expenses {
			if !expense.ExpenseDate.Before(startOfMonth) && expense.ExpenseDate.Before(endOfMonth) {
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


