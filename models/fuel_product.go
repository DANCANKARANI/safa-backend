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

//set fuel price in a station
func CreateStationFuelPrice(c *fiber.Ctx) error {
	var inputs []StationFuelProduct

	if err := c.BodyParser(&inputs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	for i := range inputs {
		if inputs[i].StationID == uuid.Nil || inputs[i].FuelProductID == uuid.Nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "station_id and fuel_product_id are required",
			})
		}
		inputs[i].ID = uuid.New()
		inputs[i].CreatedAt = time.Now()
		inputs[i].UpdatedAt = time.Now()
	}

	if err := db.Create(&inputs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create station fuel prices",
		})
	}

	return c.JSON(inputs)
}

//prices by station
func GetStationFuelPrices(c *fiber.Ctx) error {
	stationIDParam := c.Params("station_id")
	stationID, err := uuid.Parse(stationIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid station ID"})
	}

	var prices []StationFuelProduct
	if err := db.Where("station_id = ?", stationID).Preload("FuelProduct").Find(&prices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch prices",
		})
	}

	return c.JSON(prices)
}

//update prices
func UpdateStationFuelPrice(c *fiber.Ctx) error {
	idParam := c.Params("id")
	priceID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var payload struct {
		UnitPrice     float64   `json:"unit_price"`
		EffectiveFrom time.Time `json:"effective_from"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var product StationFuelProduct
	if err := db.First(&product, "id = ?", priceID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Price not found"})
	}

	product.UnitPrice = payload.UnitPrice
	product.EffectiveFrom = payload.EffectiveFrom
	product.UpdatedAt = time.Now()

	if err := db.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Update failed"})
	}

	return c.JSON(product)
}

//delete fuel prices
func DeleteStationFuelPrice(c *fiber.Ctx) error {
	idParam := c.Params("id")
	priceID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	if err := db.Delete(&StationFuelProduct{}, "id = ?", priceID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Delete failed"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

