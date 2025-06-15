package models

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

//create tank
func CreateTank(c *fiber.Ctx, tank Tank) (*Tank, error) {
	db.AutoMigrate(&Tank{})
	tank.ID = uuid.New()
	if err := db.Create(&tank).Error; err != nil {
		return nil, errors.New("failed to create tank")
	}
	return &tank, nil
}

//update tank
func UpdateTank(c *fiber.Ctx, tank Tank, id uuid.UUID) (*Tank, error) {
	if err := db.First(&tank, "id = ?", id).Error; err != nil {
		return nil, errors.New("tank not found")
	}
	if err := db.Model(&Tank{}).Where("id = ?", id).Updates(tank).Error; err != nil {
		return nil, errors.New("failed to update tank")
	}
	return &tank, nil
}

//delete tank
func DeleteTank(c *fiber.Ctx, id uuid.UUID) error {
	var tank Tank
	if err := db.First(&tank, "id = ?", id).Error; err != nil {
		return errors.New("tank not found")
	}
	if err := db.Delete(&tank).Error; err != nil {
		return errors.New("failed to delete tank")
	}
	return nil
}

//assign tank to pump
func AssignTankToPump (c *fiber.Ctx, tank_id, pump_id uuid.UUID)(*Tank, error){
	var tank Tank
	db.AutoMigrate(&Tank{})
	if err := db.First(&tank, "id = ?", tank_id).Error; err != nil {
		return nil, errors.New("tank not found")
	}
	var pump Pump
	if err := db.First(&pump, "id = ?", pump_id).Error; err != nil {
		return nil, errors.New("pump not found")
	}
	if err := db.Model(&pump).Association("Tanks").Append(&tank); err != nil {
		return nil, errors.New("failed to assign tank to pump")
	}
	return &tank, nil
}

func GetTanksByStation(c *fiber.Ctx, station_id uuid.UUID) ([]Tank, error) {
	var tanks []Tank
	err := db.Preload("Pumps").
		Preload("Dippings").
		Where("station_id = ?", station_id).
		Find(&tanks).Error

	if err != nil {
		return nil, errors.New("failed to get tanks")
	}
	return tanks, nil
}

func GetTankByID(c *fiber.Ctx, id uuid.UUID) (*Tank, error) {
	var tank Tank
	if err := db.First(&tank, "id = ?", id).Error; err != nil {
		return nil, errors.New("tank not found")
	}
	return &tank, nil
}