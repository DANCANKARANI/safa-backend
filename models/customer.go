package models

import "github.com/gofiber/fiber/v2"

// add customer
func (c *Customer) TableName() string {
	return "customers"
}
func AddCustomer(c *fiber.Ctx) error {
	var customer Customer
	// Parse the request body into the customer struct
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Create(&customer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(customer)
}

// GetCustomers retrieves all customers from the database
func GetCustomers(c *fiber.Ctx) error {
	var customers []Customer
	if err := db.Find(&customers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(customers)
}
// GetCustomer retrieves a customer by ID
func GetCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
	}
	return c.JSON(customer)
}
// UpdateCustomer updates a customer by ID
func UpdateCustomer(c *fiber.Ctx) error {
	id := c.Params("id")
	var customer Customer
	if err := db.First(&customer, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
	}
	// Parse the request body into the customer struct
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Save(&customer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(customer)
}

//record customerCredits
func (c *CustomerCredit) TableName() string {
	return "customer_credits"
}
// AddCustomerCredit adds a credit record for a customer
func AddCustomerCredit(c *fiber.Ctx) error {
	var credit CustomerCredit
	// Parse the request body into the credit struct
	if err := c.BodyParser(&credit); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Create(&credit).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(credit)
}
// GetCustomerCredits retrieves all credit records for a customer
func GetCustomerCredits(c *fiber.Ctx) error {
	var credits []CustomerCredit
	customerID := c.Params("id")
	if err := db.Where("customer_id = ?", customerID).Find(&credits).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(credits)
}

//add customer credit payments
func (c *CustomerCreditPayment) TableName() string {
	return "customer_credit_payments"
}

// AddCustomerCreditPayment adds a payment record for a customer credit
func AddCustomerCreditPayment(c *fiber.Ctx) error {
	var payment CustomerCreditPayment
	// Parse the request body into the payment struct
	if err := c.BodyParser(&payment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := db.Create(&payment).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(payment)
}

// GetCustomerCreditPayments retrieves all payment records for a customer credit
func GetCustomerCreditPayments(c *fiber.Ctx) error {
	var payments []CustomerCreditPayment
	creditID := c.Params("id")
	if err := db.Where("customer_credit_id = ?", creditID).Find(&payments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(payments)
}
// GetCustomerCreditPayment retrieves a payment record by ID
func GetCustomerCreditPayment(c *fiber.Ctx) error {
	id := c.Params("id")
	var payment CustomerCreditPayment
	if err := db.First(&payment, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Payment not found"})
	}
	return c.JSON(payment)
}