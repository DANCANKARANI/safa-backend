package models

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateNozzle (c *fiber.Ctx, nozzle Nozzle)(*Nozzle, error) {
	nozzle.ID = uuid.New()
	err := db.Create(&nozzle).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to create nozzle")
	}
	return &nozzle, nil
}

//update Nozzle
func UpdateNozzle(c *fiber.Ctx,nozzle Nozzle, id uuid.UUID) (*Nozzle, error) {

	err := db.Model(&Nozzle{}).Where("id = ?", id).Updates(&nozzle).Error
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to update nozzle")
	}
	return &nozzle, nil
}

//delete nozzle
func DeleteNozzle(id uuid.UUID)error{
	nozzle := new(Nozzle)
	if err := db.First(&nozzle, "id = ?", id).Error; err != nil {
		return errors.New("nozzle not found")
	}
	if err := db.Delete(&nozzle).Error; err != nil {
		return errors.New("failed to delete nozzle")
	}
	return nil
}

//get all nozzles
func GetAllNozzles() ([]Nozzle, error) {
	var nozzles []Nozzle
	if err := db.Preload("Sales").Find(&nozzles).Error; err != nil {
		return nil, errors.New("failed to get nozzles")
	}
	return nozzles, nil
}

//get nozzle by id
func GetNozzleByID(id uuid.UUID) (*Nozzle, error) {
	nozzle := new(Nozzle)
	if err := db.Preload("Sales").First(&nozzle, "id = ?", id).Error; err != nil {
		return nil, errors.New("nozzle not found")
	}
	return nozzle, nil
}