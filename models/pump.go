package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreatePump(c *fiber.Ctx, pump *Pump)(*Pump, error) {

	pump.ID = uuid.New()
	err := db.Create(&pump).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to create pump")
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
