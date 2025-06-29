package models

func MigrateDb(){
	// Perform database migration tasks here
	db.AutoMigrate(
		&StationFuelProduct{},
		&Employee{},
		&Payment{},
		&SalaryAdvance{},
		&Station{},
		&Expenses{},
		&Nozzle{},


		&PumpReadings{},
		&Sales{},
		&Dippings{},
		&FuelStock{},
		&Supplier{},
		
		&SupplierPayment{},
		&SupplierDebt{},
		&FuelProduct{},
		&FuelTransaction{},
		&Supply{},
		&Pump{},
		
		&Tank{},
	)
}