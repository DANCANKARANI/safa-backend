package models

import (
	"errors"
	"fmt"
	"log"
	"github.com/google/uuid"
)

func AddEmployeePayment(p *Payment) (*Payment,error) {
	p.ID = uuid.New()
	err := db.Create(p).Error
	if err != nil {
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
