package repositories

import (
	"errors"
	"fmt"

	"github.com/dancankarani/safa/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func RecordSupply(db *gorm.DB, supply models.Supply) error {
    return db.Transaction(func(tx *gorm.DB) error {
       

        // 2. Get current credit balance
        var supplier models.Supplier
        db.AutoMigrate(&supplier)
        if err := tx.First(&supplier, supply.SupplierID).Error; err != nil {
            return err
        }

        // 3. Calculate new balance (apply any credit first)
        amountDue := supply.TotalAmount
        newBalance := supplier.CreditBalance - amountDue

        var debtRecord models.SupplierDebt
        if newBalance >= 0 {
            // Fully covered by credit
            debtRecord = models.SupplierDebt{
                ID: uuid.New(),
                SupplierID:     supply.SupplierID,
                SupplyID:       supply.ID,
                TransactionType: "supply",
                Amount:        amountDue,
                RunningBalance: -newBalance, // New credit balance
                Notes:         fmt.Sprintf("Paid from credit. Remaining credit: %.2f", newBalance),
            }
            supplier.CreditBalance = newBalance
        } else {
            // Partial credit, create debt
            debtRecord = models.SupplierDebt{
                ID: uuid.New(),
                SupplierID:     supply.SupplierID,
                SupplyID:       supply.ID,
                TransactionType: "supply",
                Amount:        amountDue,
                RunningBalance: -newBalance, // Positive = debt
                Notes:         fmt.Sprintf("Applied credit: %.2f", supplier.CreditBalance),
            }
            supplier.CreditBalance = 0
        }

        // 4. Update supplier and create debt record
        if err := tx.Model(&supplier).Update("credit_balance", supplier.CreditBalance).Error; err != nil {
            return err
        }
        if err := tx.Create(&debtRecord).Error; err != nil {
            return err
        }

        // 5. Update fuel stock
        return AddFuelTransaction(tx, "supply", supply.FuelProductID, supply.StationID, supply.Quantity, supply.ID, supply.EmployeeID)
    })
}


func RecordSupplierPayment(db *gorm.DB, payment models.SupplierPayment) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Validate
        if payment.Amount <= 0 {
            return fmt.Errorf("invalid payment amount")
        }

        payment.ID = uuid.New()
        if err := tx.AutoMigrate(&models.SupplierPayment{}); err != nil {
            return err
        }

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
        tx.Model(&models.SupplierDebt{}).
            Select("COALESCE(MAX(running_balance), 0) as balance").
            Where("supplier_id = ?", payment.SupplierID).
            Scan(&currentDebt)

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
            Amount:          payment.Amount, // Store as positive for clarity
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

//get all debts
func GetSupplierDebts() ([]models.SupplierDebt, error) {
    var debts []models.SupplierDebt
    if err := db.Model(&debts).Error; err != nil {
        return nil, err
    }
    return debts, nil
}
