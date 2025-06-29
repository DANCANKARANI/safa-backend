package repositories

import (

	"fmt"

	"github.com/dancankarani/safa/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func AddPumpReadings(c *fiber.Ctx, pumpReadings models.PumpReadings) (*models.PumpReadings, error) {

	var pump models.Pump
	err := db.Preload("Tanks.FuelProduct").
		Where("id = ?", pumpReadings.PumpID).
		First(&pump).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find pump: %w", err)
	}

	if len(pump.Tanks) == 0 {
		return nil, fmt.Errorf("no tanks linked to pump %s", pumpReadings.PumpID)
	}

	// Step 1: Calculate UnitPrice

	

	// Step 2: Set UnitPrice on pumpReadings

	returnVal := &pumpReadings

	err = db.Transaction(func(tx *gorm.DB) error {
		// Auto generate ID
		pumpReadings.ID = uuid.New()

		var pump models.Pump
		if err := tx.Preload("Tanks").First(&pump, "id = ?", pumpReadings.PumpID).Error; err != nil {
			return fmt.Errorf("failed to fetch pump or related tanks: %v", err)
		}

		if len(pump.Tanks) == 0 {
			return fmt.Errorf("no tanks associated with this pump")
		}

		tank := pump.Tanks[0]

		// 🛡 Validate tank data
		if tank.ID == uuid.Nil || tank.StationID == uuid.Nil || tank.FuelProductID == uuid.Nil {
			return fmt.Errorf("incomplete tank configuration")
		}
		// 1. Save pump readings
		fmt.Println(pumpReadings.LitersDispensed)
		
		if pumpReadings.LitersDispensed < 0 {
			return fmt.Errorf("liters dispensed cannot be negative")
		}
		if err := tx.Create(&pumpReadings).Error; err != nil {
			return err
		}

		// 2. Create sales record
		sale := models.Sales{
			ID:            uuid.New(),
			EmployeeID:    pumpReadings.RecordedBy,
			PumpID:        pumpReadings.PumpID,
			LitersSold:    pumpReadings.LitersDispensed,
			PricePerLiter: pumpReadings.UnitPrice,
			TotalAmount:   pumpReadings.TotalSalesAmount,
		}

		if err := tx.Create(&sale).Error; err != nil {
			return err
		}

		// 3. Update stock using tank/station/fuel product
		if err := UpdateFuelStock(tx, tank.ID, tank.StationID, tank.FuelProductID, pumpReadings.LitersDispensed, "out"); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return returnVal, nil
}
