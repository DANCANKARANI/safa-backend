package repositories

import (
	"errors"
	"log"

	"github.com/dancankarani/safa/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


//adds new supply
func AddNewSupply(c *fiber.Ctx)(*models.Supply, error){
	var supply *models.Supply
	db.AutoMigrate(&supply)
	if err := c.BodyParser(&supply); err != nil {
		return nil, err
	}
	supply.ID = uuid.New()
	if err := db.Create(&supply).Error; err != nil {
		log.Println(err.Error())
		return nil,errors.New("failed to add supply")
	}
	
	//update the supplier debts
	err := RecordSupply(db, *supply)
	if err != nil {
		//rollback
		return nil, err
	}

	return supply, nil
}

//update supply
func UpdateSupply(c *fiber.Ctx, id uuid.UUID, supply *models.Supply)(*models.Supply, error){
	var updatedSupply models.Supply
	
	if err := db.First(&updatedSupply, "id = ?", id).Error; err != nil {
		return nil, err
	}	
	updatedSupply.TotalAmount = (updatedSupply.Quantity * updatedSupply.UnitPrice)

	log.Println(updatedSupply)

	
	updatedSupply.SupplierID = supply.SupplierID
	updatedSupply.StationID = supply.StationID
	updatedSupply.EmployeeID = supply.EmployeeID
	updatedSupply.ReferenceNo = supply.ReferenceNo
	updatedSupply.FuelProductID = supply.FuelProductID
	updatedSupply.Quantity = supply.Quantity
	updatedSupply.UnitPrice = supply.UnitPrice
	updatedSupply.DeliveryDate = supply.DeliveryDate
	updatedSupply.IsPaid = supply.IsPaid
	log.Println(updatedSupply.TotalAmount)
	
	if err := db.Updates(&updatedSupply).Scan(&updatedSupply).Error; err != nil {
		return nil, err
	}
	//record debts
	updatedSupply.ID = uuid.New()
	updatedSupply.TotalAmount = (updatedSupply.Quantity * updatedSupply.UnitPrice)
	err := RecordSupply(db, updatedSupply)
	if err != nil {
		//rollback
		return nil, err
	}
	return &updatedSupply, nil
}

//delete supply
func DeleteSupply(c *fiber.Ctx, id uuid.UUID) error {
	var supply models.Supply
	if err := db.First(&supply, "id = ?", id).Error; err != nil {
		return err
	}
	if err := db.Delete(&supply).Error; err != nil {
		return err
	}
	return nil
}

//Get supply by id
func GetSupplyByID(c *fiber.Ctx, id uuid.UUID)(*models.Supply, error){
	var supply models.Supply
	if err := db.First(&supply, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &supply, nil
}

//get all supplies
func GetAllSupplies(c *fiber.Ctx) (*[]models.Supply, error) {
	var supplies []models.Supply

	// Get pagination query params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	if err := db.Limit(limit).Offset(offset).Find(&supplies).Error; err != nil {
		return nil, err
	}
	return &supplies, nil
}

