package models

import (
	"errors"
	"log"
	"math"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (d *Dippings) BeforeSave(tx *gorm.DB) (err error) {
	d.LitersDispensed = d.OpeningDip+(d.AmountSupplied) - d.ClosingDip
	return nil
}

// Dippings represents the dipping records for fuel tanks
func (d *Dippings) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New()
	return nil
}

func CreateDipping ( dippings *Dippings)(*Dippings, error) {
	if err := db.Create(dippings).Error; err != nil {
		return nil,err
	}
	return dippings, nil
}

func GetDippingByID(id uuid.UUID) (*Dippings, error) {
	var dippings Dippings
	if err := db.First(&dippings, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dippings, nil
}

func GetDippingByStationID(id uuid.UUID) ([]Dippings, error) {
	var dippings []Dippings
	if err := db.Find(&dippings, "station_id = ?", id).Error; err != nil {
		return nil, err
	}
	return dippings, nil
}
//get dippings by fuel product
func GetDippingByFuelProductID(id uuid.UUID) ([]Dippings, error) {
	var dippings []Dippings
	if err := db.Find(&dippings, "fuel_product_id = ?", id).Error; err != nil {
		return nil, err
	}
	return dippings, nil
}
//get dippings by date
func GetDippingByDippingDate(date time.Time) ([]Dippings, error) {
	var dippings []Dippings
	if err := db.Find(&dippings, "dipping_date = ?", date).Error; err != nil {
		return nil, err
	}
	return dippings, nil
}


//delete
func DeleteDipping(c *fiber.Ctx, id uuid.UUID) error {
	var dippings Dippings
	if err := db.First(&dippings, "id = ?", id).Error; err != nil {
		return err
	}
	if err := db.Delete(&dippings).Error; err != nil {
		return err
	}
	return nil
}
 
// get all paginated dippings
type DippingResponse struct {
	Dippings
	UnitCost         float64 `json:"unit_cost"`
	AmountLostGained float64 `json:"amount_lost_gained"`
}


// func GetUnitCostAtDate(stationID, fuelProductID uuid.UUID, date time.Time) (float64, error)

func GetAllDippings(c *fiber.Ctx) ([]DippingResponse, error) {
	var dippings []Dippings

	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// Handle month/year filtering
	month := c.QueryInt("month", 0)
	year := c.QueryInt("year", 0)
	now := time.Now()
	if month < 1 || month > 12 {
		month = int(now.Month())
	}
	if year < 1 {
		year = now.Year()
	}
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 1, 0)

	// Preload Tank + Tank's FuelProduct + Tank's Station, filter by dipping_date
	if err := db.
		Preload("Tank").
		Preload("Tank.FuelProduct").
		Preload("Tank.Station").
		Where("dipping_date >= ? AND dipping_date < ?", startDate, endDate).
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&dippings).Error; err != nil {
		return nil, err
	}

	var responses []DippingResponse

	for _, dip := range dippings {
		unitCost, err := GetUnitCostAtDate(dip.Tank.StationID, dip.Tank.FuelProductID, dip.DippingDate)
		if err != nil {
			unitCost = 0
		}

		amountLostGained := dip.Deviation * unitCost

		resp := DippingResponse{
			Dippings:        dip,
			UnitCost:        unitCost,
			AmountLostGained: amountLostGained,
		}

		responses = append(responses, resp)
	}

	return responses, nil
}

//function to get deviation and AmountLostGained for every Fuel product In a munth
type MonthlyDippingSummary struct {
	FuelProductID uuid.UUID `json:"fuel_product_id"`
	FuelProduct   string    `json:"fuel_product"` // Name of the fuel product
	Deviation     float64   `json:"deviation"`
	AmountLostGained float64 `json:"amount_lost_gained"`
}

func GetMonthlyDippingSummary(c *fiber.Ctx) ([]MonthlyDippingSummary, error) {
	month := c.Query("month")
	year := c.Query("year")

	var summaries []MonthlyDippingSummary

	if err := db.
		Model(&Dippings{}).
		Select("fuel_product_id, SUM(deviation) as deviation, SUM(amount_lost_gained) as amount_lost_gained").
		Where("EXTRACT(MONTH FROM dipping_date) = ? AND EXTRACT(YEAR FROM dipping_date) = ?", month, year).
		Group("fuel_product_id").
		Find(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

type TankComparison struct {
	TankID            uuid.UUID `json:"tank_id"`
	TankName          string    `json:"tank_name"`
	StationName       string    `json:"station_name"` // ✅ Added
	ProductName       string    `json:"product_name"` // ✅ Added field
	StartDate         string    `json:"start_date"`
	EndDate           string    `json:"end_date"`
	DippingLiters     float64   `json:"dipping_liters"`
	PumpReadingLiters float64   `json:"pump_reading_liters"`
	Difference        float64   `json:"difference"`
	Matches           bool      `json:"matches"`
}
type ResTankComparison struct {
	Comparisons []TankComparison `json:"comparisons"`
}

func ComparePumpReadingsWithDippings(c *fiber.Ctx) (*ResTankComparison, error) {
	stationIDParam := c.Query("station_id")
	if stationIDParam == "" {
		return nil, errors.New("station_id query parameter is required")
	}

	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		return nil, errors.New("invalid station_id format")
	}

	var station Station
	if err := db.
		Preload("Tanks.FuelProduct").
		Preload("Tanks.Dippings").
		Preload("Tanks.Pumps.Readings").
		First(&station, "id = ?", stationID).Error; err != nil {
		return nil, errors.New("station not found")
	}

	var comparisons []TankComparison

	for _, tank := range station.Tanks {
		// Sort dippings by date
		sort.Slice(tank.Dippings, func(i, j int) bool {
			return tank.Dippings[i].DippingDate.Before(tank.Dippings[j].DippingDate)
		})

		for i := 0; i < len(tank.Dippings)-1; i++ {
			d1 := tank.Dippings[i]
			d2 := tank.Dippings[i+1]

			start := d1.DippingDate
			end := d2.DippingDate

			var totalPumpLiters float64

			for _, pump := range tank.Pumps {
				for _, reading := range pump.Readings {
					if !reading.CreatedAt.Before(start) && reading.CreatedAt.Before(end) {
						totalPumpLiters += reading.LitersDispensed
					}
				}
			}

			dippingDiff := d1.ClosingDip - d2.ClosingDip
			diff := dippingDiff - totalPumpLiters

			comparisons = append(comparisons, TankComparison{
			TankID:            tank.ID,
			TankName:          tank.Name,
			StationName:       station.Name, // ✅ Here
			ProductName:       tank.FuelProduct.Name,
			StartDate:         start.Format("2006-01-02 15:04"),
			EndDate:           end.Format("2006-01-02 15:04"),
			DippingLiters:     dippingDiff,
			PumpReadingLiters: totalPumpLiters,
			Difference:        diff,
			Matches:           math.Abs(diff) < 0.01,
		})
		}
	}

	return &ResTankComparison{Comparisons: comparisons}, nil
}

func GetUnitCostAtDate(stationID, fuelProductID uuid.UUID, date time.Time) (float64, error) {
	var sfp StationFuelProduct
	err := db.Where("station_id = ? AND fuel_product_id = ? AND effective_from <= ?", stationID, fuelProductID, date).
		Order("effective_from DESC").
		First(&sfp).Error
	if err != nil {
		return 0, err
	}
	return sfp.UnitPrice, nil
}

//get openingDipping readings
type resOpeningDippings struct {
	OpeningDip         float64   `json:"opening_dip"`
	OpeningMeterReading float64  `json:"opening_meter_reading"`
	Date               time.Time `json:"date"`
}

func GetOpeningDippings(c *fiber.Ctx, tankID uuid.UUID) (*resOpeningDippings, error) {
	var dippings Dippings
	if err := db.Where("tank_id = ?", tankID).Order("dipping_date DESC").First(&dippings).Error; err != nil {
		return nil, err
	}

	openingDippings := &resOpeningDippings{
		OpeningDip:         dippings.ClosingDip,
		OpeningMeterReading: dippings.ClosingMeter,
		Date:               dippings.DippingDate,
	}

	return openingDippings, nil
}

//get latest summation of pump readings that acts as the closing sales for the dipping


func GetLatestReadingsSumByTankID(tankID uuid.UUID) (*ResSales, error) {
	var res ResSales

	query := `
	SELECT 
		SUM(pr.closing_sales_amount) AS total_sales,
		SUM(pr.closing_meter) AS total_liters
	FROM pump_readings pr
	JOIN (
		SELECT pump_id, MAX(reading_date) AS latest_reading
		FROM pump_readings
		GROUP BY pump_id
	) latest ON pr.pump_id = latest.pump_id AND pr.reading_date = latest.latest_reading
	JOIN tank_pumps tp ON tp.pump_id = pr.pump_id
	WHERE tp.tank_id = ?
	`

	if err := db.Raw(query, tankID).Scan(&res).Error; err != nil {
		log.Println("Error fetching latest readings sum:", err)
		return nil, errors.New("failed to get latest pump readings for tank")
	}

	return &res, nil
}

func UpdateDippings(c *fiber.Ctx, dippingID uuid.UUID) (*Dippings, error) {
	var updateData Dippings
	if err := c.BodyParser(&updateData); err != nil {
		return nil, errors.New("failed to parse request body")
	}

	var existing Dippings
	if err := db.First(&existing, "id = ?", dippingID).Error; err != nil {
		return nil, errors.New("dipping not found")
	}

	// Only update fields that are non-zero or non-default
	if updateData.OpeningDip != 0 {
		existing.OpeningDip = updateData.OpeningDip
	}
	if updateData.ClosingDip != 0 {
		existing.ClosingDip = updateData.ClosingDip
	}
	if updateData.AmountSupplied != 0 {
		existing.AmountSupplied = updateData.AmountSupplied
	}
	if !updateData.DippingDate.IsZero() {
		existing.DippingDate = updateData.DippingDate
	}
	if updateData.ClosingMeter != 0 {
		existing.ClosingMeter = updateData.ClosingMeter
	}
	// Add more fields as needed

	// Recalculate LitersDispensed if relevant fields changed
	existing.LitersDispensed = existing.OpeningDip + existing.AmountSupplied - existing.ClosingDip

	if err := db.Save(&existing).Error; err != nil {
		return nil, errors.New("failed to update dipping")
	}

	return &existing, nil
}

/*delete dipping*/
