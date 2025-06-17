package repositories

import (
	"errors"
	"time"

	"github.com/dancankarani/safa/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetCurrentStock(db *gorm.DB, productID, stationID uuid.UUID) (float64, error) {
	var stock float64
	err := db.Model(&models.FuelTransaction{}).
		Select("COALESCE(SUM(CASE WHEN type = 'supply' THEN quantity ELSE -quantity END), 0)").
		Where("fuel_product_id = ? AND station_id = ?", productID, stationID).
		Scan(&stock).Error
	return stock, err
}



// direction: "in" for supply, "out" for sale or usage
func UpdateFuelStock(tx *gorm.DB, tankID, stationID, fuelProductID uuid.UUID, quantity float64, direction string) error {
	var stock models.FuelStock
	err := tx.Where("tank_id = ?", tankID).First(&stock).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if direction == "out" {
			return errors.New("cannot reduce stock: no existing stock record")
		}

		// Create new stock record
		stock = models.FuelStock{
			ID:            uuid.New(),
			TankID:        tankID,
			FuelProductID: fuelProductID,
			StationID:     stationID,
			CurrentVolume: quantity,
		}
		return tx.Create(&stock).Error
	} else if err != nil {
		return err
	}

	// Update existing stock
	if direction == "in" {
		stock.CurrentVolume += quantity
	} else if direction == "out" {
		if stock.CurrentVolume < quantity {
			return errors.New("not enough stock to deduct")
		}
		stock.CurrentVolume -= quantity
	} else {
		return errors.New("invalid stock direction")
	}

	return tx.Save(&stock).Error
}

type FuelStockWithDetails struct {
    ID             uuid.UUID `json:"id"`
    FuelProductID  uuid.UUID `json:"fuel_product_id"`
    StationID      uuid.UUID `json:"station_id"`
    TankID         uuid.UUID `json:"tank_id"`
    FuelProductName string    `json:"fuel_product_name"`
    StationName     string    `json:"station_name"`
    TankName        string    `json:"tank_name"`
    CurrentVolume   float64   `json:"current_volume"`
    LastUpdated     time.Time `json:"last_updated"`
}

func GetFuelStocksWithDetails() ([]FuelStockWithDetails, error) {
    var stocks []FuelStockWithDetails
    err := db.Table("fuel_stocks").
    Select("fuel_stocks.id, fuel_stocks.fuel_product_id, fuel_stocks.station_id, fuel_stocks.tank_id, fuel_stocks.current_volume, fuel_stocks.last_updated, "+
           "fuel_products.name as fuel_product_name, stations.name as station_name, tanks.name as tank_name").
    Joins("LEFT JOIN fuel_products ON fuel_products.id = fuel_stocks.fuel_product_id").
    Joins("LEFT JOIN stations ON stations.id = fuel_stocks.station_id").
    Joins("LEFT JOIN tanks ON tanks.id = fuel_stocks.tank_id").
    Scan(&stocks).Error

    return stocks, err
}


