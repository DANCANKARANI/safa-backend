package repositories

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/dancankarani/safa/models"
)
func AddFuelTransaction(db *gorm.DB, txType string, productID, stationID uuid.UUID, quantity float64, referenceID, userID uuid.UUID) error {
    db.AutoMigrate(&models.FuelTransaction{})
    // Calculate previous level
    previousLevel, err := GetCurrentStock(db, productID, stationID)
    if err != nil {
        return err
    }

    // Calculate new level
    var newLevel float64
    switch txType {
    case "supply":
        newLevel = previousLevel + quantity
    case "sale", "dipping":
        newLevel = previousLevel - quantity
        if newLevel < 0 {
            return fmt.Errorf("insufficient stock: only %.2f available", previousLevel)
        }
    default:
        return fmt.Errorf("invalid transaction type: %s", txType)
    }

    // Record the transaction
    transaction := models.FuelTransaction{
        ID: uuid.New(),
        FuelProductID: productID,
        StationID:    stationID,
        Type:         txType,
        Quantity:     quantity,
        PreviousLevel: previousLevel,
        NewLevel:     newLevel,
        ReferenceID:  referenceID,
        CreatedBy:    userID,
    }

    return db.Create(&transaction).Error
}

