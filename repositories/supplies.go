package repositories

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/dancankarani/safa/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RecordSupply(db *gorm.DB, supply models.Supply) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Save the supply record
		if err := tx.Create(&supply).Error; err != nil {
			return err
		}

		// 2. Load the supplier
		var supplier models.Supplier
		if err := tx.First(&supplier, "id = ?", supply.SupplierID).Error; err != nil {
			return err
		}

		// 3. Get current running balance
		var currentDebt struct {
			Balance float64
		}
		err := tx.Model(&models.SupplierDebt{}).
			Select(`COALESCE(SUM(CASE 
				WHEN transaction_type = 'supply' THEN amount 
				WHEN transaction_type = 'payment' THEN -amount 
				ELSE 0 END), 0) as balance`).
			Where("supplier_id = ?", supply.SupplierID).
			Scan(&currentDebt).Error
		if err != nil {
			return err
		}

		// 4. Apply credit
		amountDue := supply.TotalAmount
		creditToApply := math.Min(supplier.CreditBalance, amountDue)
		remainingDebt := amountDue - creditToApply
		newRunningBalance := currentDebt.Balance + remainingDebt

		// 5. Create supplier debt record
		notes := ""
		if creditToApply > 0 {
			notes = fmt.Sprintf("Applied credit: %.2f, remaining debt: %.2f", creditToApply, remainingDebt)
		} else {
			notes = "No credit applied"
		}

		debtRecord := models.SupplierDebt{
			ID:              uuid.New(),
			SupplierID:      supply.SupplierID,
			SupplyID:        &supply.ID,
			TransactionType: "supply",
			Amount:          amountDue,
			RunningBalance:  newRunningBalance,
			Notes:           notes,
		}

		// 6. Update supplier's credit balance
		supplier.CreditBalance -= creditToApply
		if supplier.CreditBalance < 0 {
			supplier.CreditBalance = 0
		}
		if err := tx.Model(&supplier).Update("credit_balance", supplier.CreditBalance).Error; err != nil {
			return err
		}

		// 7. Record the debt entry
		if err := tx.Create(&debtRecord).Error; err != nil {
			return err
		}

		// 8. Update fuel stock
		if err := UpdateFuelStock(tx, supply.TankID, supply.StationID, supply.FuelProductID, supply.Quantity, "in"); err != nil {
			return err
		}

		// 9. Add fuel transaction
		return AddFuelTransaction(tx, "supply", supply.FuelProductID, supply.StationID, supply.Quantity, supply.ID, supply.EmployeeID)
	})
}


func RecordSupplierPayment(db *gorm.DB, payment models.SupplierPayment) error {
    return db.Transaction(func(tx *gorm.DB) error {
        if payment.Amount <= 0 {
            return fmt.Errorf("invalid payment amount")
        }

        payment.ID = uuid.New()
        if err := tx.Create(&payment).Error; err != nil {
            return err
        }

        var supplier models.Supplier
        if err := tx.First(&supplier, "id = ?", payment.SupplierID).Error; err != nil {
            return err
        }

        var currentDebt struct {
            Balance float64
        }

        // Get net debt (supplies - payments)
        err := tx.Model(&models.SupplierDebt{}).
            Select(`COALESCE(SUM(CASE 
                WHEN transaction_type = 'supply' THEN amount 
                WHEN transaction_type = 'payment' THEN -amount 
                ELSE 0 END), 0) as balance`).
            Where("supplier_id = ?", payment.SupplierID).
            Scan(&currentDebt).Error
        if err != nil {
            return err
        }

        newBalance := currentDebt.Balance - payment.Amount
        notes := ""

        if newBalance < 0 {
            notes = fmt.Sprintf("Overpayment. New credit: %.2f", -newBalance)
            supplier.CreditBalance = -newBalance
            newBalance = 0
        } else {
            notes = fmt.Sprintf("Payment applied. Remaining debt: %.2f", newBalance)
            supplier.CreditBalance = 0
        }

        debtRecord := models.SupplierDebt{
            ID:              uuid.New(),
            SupplierID:      payment.SupplierID,
            TransactionType: "payment",
            Amount:          payment.Amount,
            RunningBalance:  newBalance,
            Notes:           notes,
        }

        if err := tx.Model(&supplier).Update("credit_balance", supplier.CreditBalance).Error; err != nil {
            return err
        }

        return tx.Create(&debtRecord).Error
    })
}


//get balances
func GetSupplierBalance(db *gorm.DB, supplierID uuid.UUID) (debtYouOwe float64, creditTheyOwe float64, netBalance float64, err error) {
    // Get supplier's credit balance
    var supplier models.Supplier
    if err = db.Select("credit_balance").First(&supplier, supplierID).Error; err != nil {
        return 0, 0, 0, err
    }
    creditTheyOwe = supplier.CreditBalance

    // Get current debt (positive = you owe supplier)
   var latestDebt models.SupplierDebt
    err = db.Model(&models.SupplierDebt{}).
        Where("supplier_id = ?", supplierID).
        Order("created_at DESC").
        Limit(1).
        Take(&latestDebt).Error

    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return 0, 0, 0, err
    }
    debtYouOwe = latestDebt.RunningBalance

    // Calculate net balance (positive = you owe, negative = they owe)
    netBalance = debtYouOwe - creditTheyOwe
    
    return debtYouOwe, creditTheyOwe, netBalance, nil
}



type SupplierDebtDTO struct {
	TransactionType    string
	FuelType           string
	Quantity           float64
	UnitCost           float64
	CarNumber          string
	Amount             float64
	RunningBalance     float64
	SellingPrice       float64
	Profit             float64
	Notes              string
	Date               time.Time
}

func GetSupplierDebts(db *gorm.DB) ([]SupplierDebtDTO, error) {
	var debts []models.SupplierDebt
	err := db.
		Preload("Supply").
		Preload("Supply.FuelProduct").
		Order("created_at DESC").
		Find(&debts).Error
	if err != nil {
		return nil, err
	}

	var result []SupplierDebtDTO
	for _, debt := range debts {
		dto := SupplierDebtDTO{
			TransactionType: debt.TransactionType,
			Amount:          debt.Amount,
			RunningBalance:  debt.RunningBalance,
			Notes:           debt.Notes,
			Date:            debt.CreatedAt,
		}

		if debt.Supply != nil {
			supply := debt.Supply
			dto.CarNumber = supply.CarNumber
			dto.Quantity = supply.Quantity
			dto.UnitCost = supply.UnitPrice
			dto.FuelType = supply.FuelProduct.Name

			// Get latest selling price from StationFuelProduct
			var sellingPrice float64
			var stationFuelProduct models.StationFuelProduct
			err := db.
				Where("fuel_product_id = ? AND station_id = ?", supply.FuelProductID, supply.StationID).
				Order("effective_from DESC").
				First(&stationFuelProduct).Error
			if err == nil {
				sellingPrice = stationFuelProduct.UnitPrice
			}

			dto.SellingPrice = sellingPrice
			dto.Profit = (sellingPrice - supply.UnitPrice) * supply.Quantity
		}

		result = append(result, dto)
	}

	return result, nil
}

//extract Running Balance
func GetSupplierRunningBalance(tx *gorm.DB, supplierID uuid.UUID) (float64, error) {
    var result struct {
        Balance float64
    }
    err := tx.Model(&models.SupplierDebt{}).
        Select(`COALESCE(SUM(CASE 
            WHEN transaction_type = 'supply' THEN amount 
            WHEN transaction_type = 'payment' THEN -amount 
            ELSE 0 END), 0) as balance`).
        Where("supplier_id = ?", supplierID).
        Scan(&result).Error
    return result.Balance, err
}
