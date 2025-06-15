package models

import (
	"log"
	"time"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)	

func (s *Sales) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}

func (s *Sales) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return
}

func (s *Sales) BeforeSave(tx *gorm.DB) (err error) {
	s.TotalAmount = s.LitersSold * s.PricePerLiter
	return
}

func AddNewSales(c *fiber.Ctx)(*Sales, error){
	var sales *Sales
	if err := c.BodyParser(&sales); err != nil {
		return nil, err
	}
	if err := db.Create(&sales).Error; err != nil {
		log.Println(err.Error())
		return nil,errors.New("failed to add sales")
	}
	return sales, nil
}

//update sales
func UpdateSales(c *fiber.Ctx, id uuid.UUID, updatedData *Sales) (*Sales, error) {
	var sales Sales
	if err := db.First(&sales, "id = ?", id).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("Sales not found")
	}

	if updatedData.LitersSold != 0 {
		sales.LitersSold = updatedData.LitersSold
	}

	if updatedData.PricePerLiter != 0 {
		sales.PricePerLiter = updatedData.PricePerLiter
	}

	if err := db.Save(&sales).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to update sales")
	}
	return &sales, nil
}

//delete sales
func DeleteSales(c *fiber.Ctx, id uuid.UUID) error {
	var sales Sales
	db.AutoMigrate(&sales)
	if err := db.First(&sales, "id = ?", id).Error; err != nil {
		log.Println(err.Error())
		return errors.New("Sales not found")
	}
	if err := db.Delete(&sales).Error; err != nil {
		log.Println(err.Error())
		return errors.New("failed to delete sales")
	}
	return nil
}
// get paginated sales, order by created_at desc (latest first)
func GetSales(c *fiber.Ctx) ([]Sales, error) {
	var sales []Sales

	// Get pagination params from query, with defaults
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&sales).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to get sales")
	}
	return sales, nil
}

//get sales by dates
func GetSalesByDate(c *fiber.Ctx, startDate, endDate time.Time) ([]Sales, error) {
	var sales []Sales
	if err := db.Where("created_at >= ? AND created_at <= ?", startDate, endDate).Find(&sales).Error; err != nil {
		log.Println(err.Error())
		return nil, errors.New("failed to get sales by date")
	}
	return sales, nil
}

//get weekly sales
