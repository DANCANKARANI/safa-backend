package models

import (
	"errors"
	"math"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (d *Dippings) BeforeSave(tx *gorm.DB) (err error) {
	d.LitersDispensed = d.OpeningDip -d.ClosingDip 
	return nil
}

func (d *Dippings) AfterFind(tx *gorm.DB) (err error) {
	d.LitersDispensed = d.ClosingDip - d.OpeningDip
	return nil
}

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
func UpdateDippings(id uuid.UUID, updated *Dippings) (*Dippings, error) {
	var dippings Dippings
	if err := db.First(&dippings, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Update fields
	if updated.OpeningDip != 0 {
		dippings.OpeningDip = updated.OpeningDip
	}
	
	if updated.ClosingDip != 0 {
		dippings.ClosingDip = updated.ClosingDip
	}
	dippings.DippingDate = updated.DippingDate
	
	// Add other fields as needed

	if err := db.Save(&dippings).Error; err != nil {
		return nil, err
	}
	return &dippings, nil
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
func GetAllDippings(c *fiber.Ctx) ([]Dippings, error) {
	var dippings []Dippings

	// Get page and pageSize from query params, set defaults if not provided
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	if err := db.Limit(pageSize).Offset(offset).Find(&dippings).Error; err != nil {
		return nil, err
	}
	return dippings, nil
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
