package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/dancankarani/safa/models"
)

func GetCurrentStock(db *gorm.DB, productID, stationID uuid.UUID) (float64, error) {
	var stock float64
	err := db.Model(&models.FuelTransaction{}).
		Select("COALESCE(SUM(CASE WHEN type = 'supply' THEN quantity ELSE -quantity END), 0)").
		Where("fuel_product_id = ? AND station_id = ?", productID, stationID).
		Scan(&stock).Error
	return stock, err
}

