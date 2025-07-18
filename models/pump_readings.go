package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (p *PumpReadings) BeforeSave(tx *gorm.DB) (err error) { 
	p.LitersDispensed = p.ClosingMeter - p.OpeningMeter
	p.TotalSalesAmount =   p.LitersDispensed * p.UnitPrice
	return nil
}



func UpdatePumpReadings(c *fiber.Ctx,id uuid.UUID, updatedReadings PumpReadings)(*PumpReadings, error) {
	pumpReadings := PumpReadings{}
	err := db.First(&pumpReadings, "id = ?", id).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("pump readings not found")
	}
	
	if !updatedReadings.BusinessDay.IsZero() {
		pumpReadings.BusinessDay = updatedReadings.BusinessDay
		// Optional: disallow future business dates
		if pumpReadings.BusinessDay.After(time.Now()) {
			return nil, fmt.Errorf("business_day cannot be in the future")
		}

		pumpReadings.BusinessDay = updatedReadings.BusinessDay
	}

	if updatedReadings.OpeningSalesAmount != 0 {
		pumpReadings.OpeningSalesAmount = updatedReadings.OpeningSalesAmount
	}
	if updatedReadings.ClosingSalesAmount != 0 {
		pumpReadings.ClosingSalesAmount = updatedReadings.ClosingSalesAmount
	}
	if updatedReadings.OpeningMeter != 0 {
		pumpReadings.OpeningMeter = updatedReadings.OpeningMeter
	}
	if updatedReadings.ClosingMeter != 0 {
		pumpReadings.ClosingMeter = updatedReadings.ClosingMeter
	}
	
	if !updatedReadings.ReadingDate.IsZero() {
		pumpReadings.ReadingDate = updatedReadings.ReadingDate
	}
	if updatedReadings.Shift != ""{
		pumpReadings.ReadingDate = updatedReadings.ReadingDate
	}
	
	err = db.Save(&pumpReadings).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to update pump readings")
	}
	return &pumpReadings, nil
}

//get latest pump readings BY STATION ID
func GetLatestPumpReadingsByStationID(c *fiber.Ctx, stationID uuid.UUID) ([]PumpReadings, error) {
	var readings []PumpReadings

	err := db.
		Joins("JOIN pumps ON pumps.id = pump_readings.pump_id").
		Where("pumps.station_id = ?", stationID).
		Order("pump_readings.created_at DESC").
		Limit(10).
		Find(&readings).Error

	if err != nil {
		return nil, errors.New("failed to get latest pump readings")
	}

	return readings, nil
}

// get paginated readings, ordered by time
type PumpReadingResponse struct {
	ID                 uuid.UUID `json:"id"`
	PumpID            uuid.UUID `json:"pump_id"`
	ReadingDate        time.Time `json:"reading_date"`
	StationName        string    `json:"station"`
	FuelType           string    `json:"fuel_type"`
	Shift              string    `json:"shift"`
	OpeningMeter       float64   `json:"opening_meter"`
	ClosingMeter       float64   `json:"closing_meter"`
	OpeningSalesAmount float64   `json:"opening_sales"`
	ClosingSalesAmount float64   `json:"closing_sales"`
	LitersDispensed    float64   `json:"liters_dispensed"`
	UnitPrice          float64   `json:"unit_price"`
	TotalSalesAmount   float64   `json:"total_sales"`
}

func GetPaginatedPumpReadings(c *fiber.Ctx, page, pageSize int) ([]PumpReadingResponse, int64, error) {
	var pumpReadings []PumpReadings
	var total int64

	// Get "month" param (format: YYYY-MM)
	monthStr := c.Query("month") // Example: "2025-05"

	// Default: current year and month
	var startOfMonth, endOfMonth time.Time
	var err error

	if monthStr != "" {
		// Parse from query
		startOfMonth, err = time.Parse("2006-01", monthStr)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid month format, use YYYY-MM")
		}
	} else {
		now := time.Now()
		startOfMonth = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	endOfMonth = startOfMonth.AddDate(0, 1, 0)

	offset := (page - 1) * pageSize

	// Count total filtered records
	if err := db.Model(&PumpReadings{}).
		Where("reading_date >= ? AND reading_date < ?", startOfMonth, endOfMonth).
		Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count pump readings")
	}

	// Fetch paginated, filtered, and preloaded records
	if err := db.
		Preload("Pump.Tanks.FuelProduct").
		Preload("Pump.Tanks.Station").
		Where("reading_date >= ? AND reading_date < ?", startOfMonth, endOfMonth).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&pumpReadings).Error; err != nil {
		return nil, 0, errors.New("failed to get paginated pump readings")
	}

	var results []PumpReadingResponse
	for _, reading := range pumpReadings {
		var stationName, fuelType string
		if len(reading.Pump.Tanks) > 0 {
			tank := reading.Pump.Tanks[0]
			stationName = tank.Station.Name
			fuelType = tank.FuelProduct.Name
		}

		results = append(results, PumpReadingResponse{
			ID:                 reading.ID,
			PumpID:            reading.PumpID,
			ReadingDate:        reading.ReadingDate,
			StationName:        stationName,
			FuelType:           fuelType,
			Shift:              reading.Shift,
			OpeningMeter:       reading.OpeningMeter,
			ClosingMeter:       reading.ClosingMeter,
			OpeningSalesAmount: reading.OpeningSalesAmount,
			ClosingSalesAmount: reading.ClosingSalesAmount,
			LitersDispensed:    reading.LitersDispensed,
			UnitPrice:          reading.UnitPrice,
			TotalSalesAmount:   reading.TotalSalesAmount,
		})
	}

	return results, total, nil
}


func DeletePumpReadings(c *fiber.Ctx, id uuid.UUID) error {
	var pumpReadings PumpReadings
	if err := db.First(&pumpReadings, "id = ?", id).Error; err != nil {
		log.Println(err.Error())
		return errors.New("pump readings not found")
	}
	if err := db.Delete(&pumpReadings).Error; err != nil {
		log.Println(err.Error())
		return errors.New("failed to delete pump readings")
	}
	return nil
}
// get total sales for a date range
type ResSales struct {
	TotalSales  float64 `json:"total_sales"`
	TotalLiters float64	`json:"total_liters"`
}

func GetTotalSalesByDate(c *fiber.Ctx) (*ResSales, error) {
	// Parse date query params
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	// If dates not provided, use today's range
	if startDateStr == "" || endDateStr == "" {
		now := time.Now()
		
		// startDate: today at 8:00 AM
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
		
		// endDate: tomorrow at 8:00 AM
		endDate = startDate.Add(24 * time.Hour)
	} else {
		// parse startDate from user input (expected YYYY-MM-DD)
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, errors.New("invalid start_date format, expected YYYY-MM-DD")
		}
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 8, 0, 0, 0, startDate.Location())
		
		// parse endDate from user input (expected YYYY-MM-DD)
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, errors.New("invalid end_date format, expected YYYY-MM-DD")
		}
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 8, 0, 0, 0, endDate.Location()).Add(24 * time.Hour)
	}


	var res ResSales
	if err := db.Model(&PumpReadings{}).
		Where("reading_date BETWEEN ? AND ?", startDate, endDate).
		Select(`
			COALESCE(SUM(total_sales_amount), 0) as total_sales, 
			COALESCE(SUM(liters_dispensed), 0) as total_liters
		`).
		Scan(&res).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to get total sales by date")
	}

	return &res, nil
}

//get paginate pump readings per station
func GetPumpReadingsByStation(c *fiber.Ctx, stationID uuid.UUID) ([]PumpReadings, error) {
	var pumpReadings []PumpReadings
	if err := db.Where("station_id = ?", stationID).Find(&pumpReadings).Error; err != nil {
		return nil, errors.New("failed to get pump readings by station")
	}
	return pumpReadings, nil
}

//get latest Meter reading and sales readings
//this will act as opening sales and meter readings for the next
type resOpeningReadings struct {
	OpeningMeter       float64 `json:"opening_meter"`
	OpeningSalesAmount float64 `json:"opening_sales_amount"`
}

func GetOpeningReadings(c *fiber.Ctx, pumpID uuid.UUID) (*resOpeningReadings, error) {
	var pumpReadings PumpReadings
	if err := db.Where("pump_id = ?", pumpID).
		Order("created_at DESC").	
		First(&pumpReadings).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to get latest pump readings")
	}	
	return &resOpeningReadings{
		OpeningMeter:       pumpReadings.ClosingMeter,
		OpeningSalesAmount: pumpReadings.ClosingSalesAmount,
	}, nil
}

//get daily sales per station
