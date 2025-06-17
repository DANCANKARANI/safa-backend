package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)



func AddEmployeePayment(c *fiber.Ctx,p *Payment) (*Payment, error) {
	// Ensure ID is set
	p.ID = uuid.New()
	p.PaidMonth = c.Query("month")
	// Automatically set status to "paid"
	p.Status = "paid"

	// Validate paid_month format (e.g., "2025-06")
	if _, err := time.Parse("2006-01", p.PaidMonth); err != nil {
		return nil, fmt.Errorf("invalid paid_month format; expected YYYY-MM")
	}

	// Check if payment for this employee and month already exists
	var existing Payment
	if err := db.Where("employee_id = ? AND paid_month = ?", p.EmployeeID, p.PaidMonth).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("payment already exists for employee %s in %s", p.EmployeeID, p.PaidMonth)
	}

	// Set payment date if not provided
	if p.PaymentDate.IsZero() {
		p.PaymentDate = time.Now()
	}

	// Save payment
	if err := db.Create(p).Error; err != nil {
		log.Printf("Error creating employee payment: %v", err)
		return nil, fmt.Errorf("error creating employee payment: %v", err)
	}

	return p, nil
}

//update employee payment
func UpdateEmployeePayment(employeePayment *Payment, id uuid.UUID) (*Payment, error) {
	err := db.First(&employeePayment, "id = ?", id).Error
	if err != nil {
		log.Printf("Error finding employee payment: %v", err)
		return nil, fmt.Errorf("error finding employee payment: %v", err)
	}
	err = db.Save(&employeePayment).Error
	if err != nil {
		log.Printf("Error updating employee payment: %v", err)
		return nil, fmt.Errorf("error updating employee payment: %v", err)
	}
	return employeePayment, nil
}

//delete the payment 
func DeletePayment(id uuid.UUID) error {
	var payment Payment
	if err := db.First(&payment, "id = ?", id).Error; err != nil {
		return errors.New("payment not found")
	}
	if err := db.Delete(&payment).Error; err != nil {
		return errors.New("failed to delete payment")
	}
	return nil
}

func GetRecentPayments(limit int) ([]Payment, error) {
	var payments []Payment
	err := db.Order("created_at desc").Limit(limit).Find(&payments).Error
	if err != nil {
		log.Printf("Error getting recent payments: %v", err)
		return nil, fmt.Errorf("error getting recent payments: %v", err)
	}
	return payments, nil
}

//get payments report and payroll
type ResPayments struct {
	EmployeeId	uuid.UUID	`json:"payment_id"`
	Name        string  `json:"name"`
	Station     string  `json:"station"`
	BasicSalary float64 `json:"basic_salary"`
	Allowances  float64 `json:"allowances"`
	Deductions  float64 `json:"deductions"`
	NetPay      float64 `json:"net_pay"`
	Status      string  `json:"status"`
	PaidMonth   string  `json:"paid_month"`
}

type ResReport struct {
	Payments          []ResPayments `json:"payments"`
	TotalPaidAmount   float64       `json:"total_paid_amount"`
	TotalUnpaidAmount float64       `json:"total_unpaid_amount"`
}
func GetReport(month string) (ResReport, error) {
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	var employees []Employee
	err := db.Preload("Station").
		Preload("Payments", "paid_month = ?", month).
		Preload("SalaryAdvance").
		Find(&employees).Error
	if err != nil {
		return ResReport{}, fmt.Errorf("failed to load payroll data: %v", err)
	}

	var results []ResPayments
	var totalPaid, totalUnpaid float64

	for _, emp := range employees {
		// Calculate deductions from advances in the given month
		var deductions float64
		for _, adv := range emp.SalaryAdvance {
			if adv.DateRequested.Format("2006-01") == month {
				deductions += adv.Amount
			}
		}

		allowances := 0.0 // Update if you support allowances
		netPay := emp.Salary + allowances - deductions

		status := "unpaid"
		var paidAmount float64

		if len(emp.Payments) > 0 {
			p := emp.Payments[0]
			status = p.Status
			if status == "paid" {
				paidAmount = p.Amount
				totalPaid += paidAmount
			}
		} else {
			totalUnpaid += netPay
		}

		results = append(results, ResPayments{
			EmployeeId:emp.ID ,
			Name:        emp.FirstName + " " + emp.LastName,
			Station:     emp.Station.Name,
			BasicSalary: emp.Salary,
			Allowances:  allowances,
			Deductions:  deductions,
			NetPay:      netPay,
			Status:      status,
			PaidMonth:   month,
		})
	}

	return ResReport{
		Payments:          results,
		TotalPaidAmount:   totalPaid,
		TotalUnpaidAmount: totalUnpaid,
	}, nil
}
