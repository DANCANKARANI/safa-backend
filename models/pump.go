package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreatePump(c *fiber.Ctx, pump *Pump, tankID uuid.UUID) (*Pump, error) {
	// Begin transaction
	tx := db.Begin()
	if tx.Error != nil {
		return nil, errors.New("failed to start transaction")
	}

	// Generate new UUID
	pump.ID = uuid.New()

	// Create the pump
	if err := tx.Create(pump).Error; err != nil {
		tx.Rollback()
		log.Println("Error creating pump:", err)
		return nil, errors.New("failed to create pump")
	}

	// Find the tank
	var tank Tank
	if err := tx.Preload("Pumps").First(&tank, "id = ?", tankID).Error; err != nil {
		tx.Rollback()
		log.Println("Tank not found:", err)
		return nil, fmt.Errorf("tank not found: %w", err)
	}

	// Associate the pump with the tank
	if err := tx.Model(&tank).Association("Pumps").Append(pump); err != nil {
		tx.Rollback()
		log.Println("Failed to assign pump to tank:", err)
		return nil, fmt.Errorf("failed to assign pump to tank: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Println("Transaction commit failed:", err)
		return nil, errors.New("failed to commit transaction")
	}

	return pump, nil
}


//update Pump
func UpdatePump(c *fiber.Ctx, id uuid.UUID) (*Pump, error) {
	db.AutoMigrate(&Pump{})
	pump := new(Pump)
	if err := c.BodyParser(&pump); err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to update pump")
	}
	err := db.Model(&Pump{}).Where("id = ?", id).Updates(pump).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to update pump")
	}
	return pump, nil
}

//delete pump
func DeletePump(c *fiber.Ctx, id uuid.UUID) error {
	db.AutoMigrate(&Pump{})
	var pump Pump
	if err := db.First(&pump, "id = ?", id).Error; err != nil {
		log.Println(err.Error())
		return errors.New("pump not found")
	}
	if err := db.Delete(&pump).Error; err != nil {
		log.Println(err.Error())
		return errors.New("failed to delete pump")
	}
	return nil
}
func GetPumpsByStation(c *fiber.Ctx, stationID uuid.UUID) (*[]Pump, error) {
	var pumps []Pump
	if err := db.
		Preload("Nozzles").
		Preload("Tanks").
		Where("station_id = ?", stationID).
		Find(&pumps).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to get pumps")
	}
	return &pumps, nil
}


//GetPumpByID
func GetPumpByID(c *fiber.Ctx, id uuid.UUID) (*Pump, error) {
	var pump Pump
	if err := db.Preload("Nozzles").First(&pump, "id = ?", id).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("pump not found")
	}
	return &pump, nil
}
//assign pump to tank
func AssignPumpToTank(c *fiber.Ctx, tankID uuid.UUID, pumpID uuid.UUID) error {
	var tank Tank
	db.AutoMigrate(&Tank{})
		if err := db.Preload("Pumps").First(&tank, "id = ?", tankID).Error; err != nil {
		return fmt.Errorf("tank not found: %v", err)
	}
	var pump Pump
	if err := db.First(&pump, "id = ?", pumpID).Error; err != nil {
		return errors.New("pump not found")
	}
	if err := db.Model(&tank).Association("Pumps").Append(&pump); err != nil {
		return fmt.Errorf("failed to assign pump to tank: %v", err)
	}
	
	return nil
}

func UnassignPumpFromTank(c *fiber.Ctx, tankID uuid.UUID, pumpID uuid.UUID) error {
	var tank Tank
	if err := db.Preload("Pumps").First(&tank, "id = ?", tankID).Error; err != nil {
		return fmt.Errorf("tank not found: %v", err)
	}

	var pump Pump
	if err := db.First(&pump, "id = ?", pumpID).Error; err != nil {
		return fmt.Errorf("pump not found: %v", err)
	}

	// Remove the association
	if err := db.Model(&tank).Association("Pumps").Delete(&pump); err != nil {
		return fmt.Errorf("failed to unassign pump from tank: %v", err)
	}

	return nil
}

type PumpInput struct {
	Name      string    `json:"name"`
	StationID uuid.UUID `json:"station_id"`
}

type CreateTankWithPumpsInput struct {
	TankName      string      `json:"name"`
	Capacity      float64     `json:"capacity"`
	StationID     uuid.UUID   `json:"station_id"`
	FuelProductID uuid.UUID   `json:"fuel_product_id"`
	Pumps         []PumpInput `json:"pumps"`
}

func CreateTankWithPumps(c *fiber.Ctx) error {
	var input CreateTankWithPumpsInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	tank := Tank{
		ID:            uuid.New(),
		Name:          input.TankName,
		Capacity:      input.Capacity,
		StationID:     input.StationID,
		FuelProductID: input.FuelProductID,
	}

	if err := db.Create(&tank).Error; err != nil {
		log.Println("Failed to create tank:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create tank"})
	}

	var createdPumps []Pump
	for _, p := range input.Pumps {
		newPump := Pump{
			ID:        uuid.New(),
			Name:      p.Name,
			StationID: p.StationID,
		}
		if err := db.Create(&newPump).Error; err != nil {
			log.Println("Failed to create pump:", err)
			continue
		}
		createdPumps = append(createdPumps, newPump)
	}

	if err := db.Model(&tank).Association("Pumps").Append(&createdPumps); err != nil {
		log.Println("Failed to assign pumps:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to assign pumps"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tank and pumps created successfully",
		"tank":    tank,
		"pumps":   createdPumps,
	})
}
