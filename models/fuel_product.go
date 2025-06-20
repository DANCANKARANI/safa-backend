package models

import (
	"errors"
	"time"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (g *FuelProduct) BeforeCreate(tx *gorm.DB) (err error) {
	g.ID = uuid.New()
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
	return
}

// CreateFuelProduct creates a new FuelProduct
func CreateFuelProduct(c *fiber.Ctx, product *FuelProduct) (*FuelProduct, error) {
	db.AutoMigrate(&FuelProduct{})
	product.ID = uuid.New()
	err := db.Create(&product).Error
	if err != nil {
		log.Println("failed to create Fuelproduct:", err.Error())
		return nil, errors.New("failed to create Fuelproduct")
	}
	return product, nil
}

// GetAllFuelProducts retrieves all FuelProducts
func GetAllFuelProducts(c *fiber.Ctx) ([]FuelProduct, error) {
	var products []FuelProduct
	if err := db.Find(&products).Error; err != nil {
		return nil, errors.New("failed to retrieve fuel products")
	}
	return products, nil
}

// GetFuelProductByID retrieves a FuelProduct by ID
func GetFuelProductByID(c *fiber.Ctx, id uuid.UUID) (*FuelProduct, error) {
	var product FuelProduct
	if err := db.First(&product, "id = ?", id).Error; err != nil {
		return nil, errors.New("fuel product not found")
	}
	return &product, nil
}

// UpdateFuelProduct updates a FuelProduct
func UpdateFuelProduct(c *fiber.Ctx, id uuid.UUID, updatedData *FuelProduct) (*FuelProduct, error) {
	var product FuelProduct
	if err := db.First(&product, "id = ?", id).Error; err != nil {
		return nil, errors.New("fuel product not found")
	}
	if updatedData.Name != "" {
		product.Name = updatedData.Name
	}
	if updatedData.Description != "" {
		product.Description = updatedData.Description
	}
	if updatedData.UnitPrice != 0 {
		product.UnitPrice = updatedData.UnitPrice
	}
	
	if err := db.Save(&product).Error; err != nil {
		return nil, errors.New("failed to update fuel product")
	}
	return &product, nil
}

// DeleteFuelProduct deletes a FuelProduct
func DeleteFuelProduct(c *fiber.Ctx, id uuid.UUID) error {
	var product FuelProduct
	if err := db.First(&product, "id = ?", id).Error; err != nil {
		return errors.New("fuel product not found")
	}
	if err := db.Delete(&product).Error; err != nil {
		return errors.New("failed to delete fuel product")
	}
	return nil
}