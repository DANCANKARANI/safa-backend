package models

func MigrateDb(){
	// Perform database migration tasks here
	db.AutoMigrate(
		&Dippings{},
		&FuelStock{},
		&Employee{},
		&Payment{},
		&Supplier{},
		&SalaryAdvance{},
		&SupplierPayment{},
		&SupplierDebt{},
		&Station{},
		&Expenses{},
		&FuelProduct{},
		&PumpReadings{},
		&FuelTransaction{},
		
		&Sales{},
		&Supply{},
		&Pump{},
		&Nozzle{},
		&Tank{},
	)
}