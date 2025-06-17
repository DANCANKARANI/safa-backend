package repositories

import (
	"errors"
	"log"

	"github.com/dancankarani/safa/database"
	"github.com/dancankarani/safa/models"
)
var db = database.ConnectDB()

func GetEmployeeById(id int) (models.Employee, error) {
	var employee models.Employee
	if err := db.Preload("Station").
		Preload("Payment").
		Preload("SalaryAdvance").
		First(&employee, id).Error; err != nil {
		return models.Employee{}, err
	}
	return employee, nil
}
//get employee by email
func GetEmployeeByEmail(email string) (models.Employee, error) {
	var employee models.Employee
	if err := db.Where("email = ?", email).First(&employee).Error; err != nil {
		return models.Employee{}, err
	}
	return employee, nil
}

//check if employee exist
func EmployeeExists(email string) (bool, error) {
	var count int64
	if err := db.Model(&models.Employee{}).Where("email = ?", email).Count(&count).Error; err != nil {
		log.Println("error checking user Existence:", err.Error())
		return false, err
	}
	if count == 0 {
		return false, errors.New("employee does not exist")
	}
	return count > 0, nil
}