package models

func MigrateDb(){
	// Perform database migration tasks here
	db.AutoMigrate(
		&Dippings{},
		&Employee{},
		&Supplier{},
		&SalaryAdvance{},
		&SupplierPayment{},
		&SupplierDebt{},
		&Station{},
		&Expenses{},
		&FuelProduct{},
		&PumpReadings{},
		&FuelTransaction{},
		&Payment{},
		&Sales{},
		&Supply{},
		&Pump{},
		&Nozzle{},
		&Tank{},
	)
}