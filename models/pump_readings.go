package models

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (p *PumpReadings) BeforeSave(tx *gorm.DB) (err error) { 
	p.LitersDispensed = p.OpeningMeter - p.ClosingMeter
	p.TotalSalesAmount =   p.LitersDispensed * p.UnitPrice
	p.BankDeposit = p.TotalSalesAmount- p.MpesaAmount
	return nil
}








func UpdatePumpReadings(c *fiber.Ctx,id uuid.UUID, updatedReadings PumpReadings)(*PumpReadings, error) {
	pumpReadings := PumpReadings{}
	err := db.First(&pumpReadings, "id = ?", id).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("pump readings not found")
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
	var pumpReadings []PumpReadings
	if err := db.Where("station_id = ?", stationID).Order("created_at desc").Limit(1).Find(&pumpReadings).Error; err != nil {
		return nil, errors.New("failed to get latest pump readings")
	}
	return pumpReadings, nil
}
// get paginated readings, ordered by time
func GetPaginatedPumpReadings(c *fiber.Ctx, page, pageSize int) ([]PumpReadings, int64, error) {
	var pumpReadings []PumpReadings
	var total int64

	// Count total records
	if err := db.Model(&PumpReadings{}).Count(&total).Error; err != nil {
		return nil, 0, errors.New("failed to count pump readings")
	}

	offset := (page - 1) * pageSize
	if err := db.Order("created_at desc").Limit(pageSize).Offset(offset).Find(&pumpReadings).Error; err != nil {
		return nil, 0, errors.New("failed to get paginated pump readings")
	}
	return pumpReadings, total, nil
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
	MpesaAmount float64	`json:"mpesa_amount"`
	BankDeposit float64	`json:"bank_deposit"`
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
			COALESCE(SUM(liters_dispensed), 0) as total_liters, 
			COALESCE(SUM(mpesa_amount), 0) as mpesa_amount,
			COALESCE(SUM(bank_deposit), 0) as bank_deposit
		`).
		Scan(&res).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to get total sales by date")
	}

	return &res, nil
}

//get total sales per station
