package models

import (
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddSupplier(c *fiber.Ctx, s *Supplier) (*Supplier, error) {
	s.ID = uuid.New()
	db.AutoMigrate(&Supplier{})
	if err := db.Create(s).Error; err != nil {
		return nil, errors.New("failed to create supplier")
	}
	return s, nil
}

//update supplier
func UpdateSupplier(c *fiber.Ctx, id uuid.UUID, updatedData *Supplier) (*Supplier, error) {
	var supplier Supplier
	if err := db.First(&supplier, "id = ?", id).Error; err != nil {
		return nil, errors.New("Supplier not found")
	}

	supplier.Name = updatedData.Name
	supplier.PhoneNumber = updatedData.PhoneNumber
	supplier.Address = updatedData.Address
	supplier.ContactName = updatedData.ContactName
	supplier.Email = updatedData.Email
	if err := db.Save(&supplier).Error; err != nil {
		return nil, errors.New("failed to update supplier")
	}
	return &supplier, nil
}
func GetSupplierByID(c *fiber.Ctx, id uuid.UUID) (*Supplier, error) {
	var supplier Supplier
	if err := db.Preload("SupplierDebts").First(&supplier, "id = ?", id).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("Supplier not found")
	}
	return &supplier, nil
}

func GetAllSuppliers(c *fiber.Ctx) ([]Supplier, error) {
	var suppliers []Supplier
	if err := db.Find(&suppliers).Error; err != nil {
		return nil, errors.New("failed to get suppliers")
	}
	return suppliers, nil
}


func DeleteSupplier(c *fiber.Ctx, id uuid.UUID) error {
	var supplier Supplier
	db.AutoMigrate(&Supplier{})
	if err := db.First(&supplier, "id = ?", id).Error; err != nil {
		log.Println(err.Error())
		return errors.New("Supplier not found")
	}
	if err := db.Delete(&supplier).Error; err != nil {
		return errors.New("failed to delete supplier")
	}
	return nil
}

