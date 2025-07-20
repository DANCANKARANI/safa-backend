package models

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddDailyAccounts(c *fiber.Ctx) error {
	var dailyAccounts DailyAccounts
	// Parse the request body into the dailyAccounts struct
	if err := c.BodyParser(&dailyAccounts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Load Nairobi location once here
	loc, err := time.LoadLocation("Africa/Nairobi")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load location"})
	}

	// Normalize BusinessDay to midnight Nairobi time
	if !dailyAccounts.BusinessDay.IsZero() {
		dailyAccounts.BusinessDay = time.Date(
			dailyAccounts.BusinessDay.Year(),
			dailyAccounts.BusinessDay.Month(),
			dailyAccounts.BusinessDay.Day(),
			0, 0, 0, 0,
			loc,
		)
	}

	// Validate required fields
	if dailyAccounts.StationID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "station_id is required"})
	}
	if dailyAccounts.BusinessDay.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "business_day is required"})
	}
	// Ensure business_day is not set to future dates
	now := time.Now().In(loc)
	if dailyAccounts.BusinessDay.After(now) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "business_day cannot be in the future"})
	}

	dailyAccounts.ID = uuid.New()

	if err := db.Create(&dailyAccounts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(dailyAccounts)
}


// GetDailyAccounts retrieves all daily accounts from the database
func GetDailyAccounts(c *fiber.Ctx) error {
	const DateFormat = "2006-01-02"
	loc, err := time.LoadLocation("Africa/Nairobi") // Adjust to your business timezone
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load location",
		})
	}

	dateParam := c.Query("date")

	var startOfDay, endOfDay time.Time
	if dateParam != "" {
		parsedDate, err := time.ParseInLocation(DateFormat, dateParam, loc)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Expected YYYY-MM-DD",
			})
		}
		startOfDay = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, loc)
		endOfDay = startOfDay.Add(24 * time.Hour)
	} else {
		now := time.Now().In(loc)
		yesterday := now.AddDate(0, 0, -1)
		startOfDay = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, loc)
		endOfDay = startOfDay.Add(24 * time.Hour)
	}

	var dailyAccounts []DailyAccounts
	if err := db.Preload("Station").Where("business_day >= ? AND business_day < ?", startOfDay, endOfDay).
		Find(&dailyAccounts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve daily accounts",
		})
	}

	return c.JSON(dailyAccounts)
}

func GetMonthlyDailyAccounts(c *fiber.Ctx) error {
	const DateFormatYYYYMM = "2006-01"
	loc, err := time.LoadLocation("Africa/Nairobi") // adjust timezone
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to load location",
		})
	}

	yymmParam := c.Query("yymm")

	var parsedTime time.Time
	if yymmParam == "" {
		// Use current month if not provided
		now := time.Now().In(loc)
		parsedTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	} else if len(yymmParam) == 7 {
		parsedTime, err = time.ParseInLocation(DateFormatYYYYMM, yymmParam, loc)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid yymm format. Use YYYY-MM format, e.g. 2025-06 for June 2025",
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid yymm parameter. Use YYYY-MM format, e.g. 2025-06 for June 2025",
		})
	}

	year := parsedTime.Year()
	month := int(parsedTime.Month())

	// Calculate start and end of the month
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
	endDate := startDate.AddDate(0, 1, 0)

	// Define a struct to hold sums
	type MonthlySums struct {
		TotalExpenses    float64 `json:"total_expenses"`
		TotalSalesAmount float64 `json:"total_sales_amount"`
		TotalBank        float64 `json:"total_bank"`
		TotalMpesa       float64 `json:"total_mpesa"`
		TotalDebtPaid    float64 `json:"total_debt_paid"`
		TotalDebtTaken   float64 `json:"total_debt_taken"`
	}

	var sums MonthlySums

	// Use raw SQL or GORM's Select for aggregation
	err = db.Model(&DailyAccounts{}).
		Select(
			"COALESCE(SUM(total_expenses),0) as total_expenses, " +
				"COALESCE(SUM(total_sales_amount),0) as total_sales_amount, " +
				"COALESCE(SUM(bank),0) as total_bank, " +
				"COALESCE(SUM(mpesa),0) as total_mpesa, " +
				"COALESCE(SUM(debt_paid),0) as total_debt_paid, " +
				"COALESCE(SUM(debt_taken),0) as total_debt_taken").
		Where("business_day >= ? AND business_day < ?", startDate, endDate).
		Scan(&sums).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve monthly sums",
		})
	}

	return c.JSON(sums)
}


func GetDailyAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	var dailyAccounts DailyAccounts
	if err := db.First(&dailyAccounts, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Daily account not found"})
	}
	return c.JSON(dailyAccounts)
}

func UpdateDailyAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	var dailyAccounts DailyAccounts
	if err := db.First(&dailyAccounts, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Daily account not found"})
	}
	// Parse the request body into the dailyAccounts struct
	if err := c.BodyParser(&dailyAccounts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Validate required fields
	if dailyAccounts.StationID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "station_id is required"})
	}
	
	// Ensure business_day is not set to future dates
	
	if !dailyAccounts.BusinessDay.IsZero() {
		// Validate business_day format
		if dailyAccounts.BusinessDay.After(time.Now()) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "business_day cannot be in the future"})
		}
		dailyAccounts.BusinessDay = dailyAccounts.BusinessDay.Truncate(24 * time.Hour) // Normalize to start of the day
	}
	dailyAccounts.ID = uuid.MustParse(id)
	// Update the dailyAccounts record in the database
	if err := db.Save(&dailyAccounts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(dailyAccounts)
}
