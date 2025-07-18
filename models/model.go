package models

import (
	"time"

	"github.com/dancankarani/safa/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var db = database.ConnectDB()
type Employee struct {
	ID        uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	FirstName string	`json:"first_name" gorm:"size:100"`
	LastName  string	`json:"last_name" gorm:"size:100"`
	Position  string	`json:"position" gorm:"size:100"`
	PhoneNumber string	`json:"phone_number" gorm:"size:15"`
	Email     string	`json:"email" gorm:"size:100"`
	StationID uuid.UUID `json:"station_id" gorm:"type:varchar(36);not null"`
	CanLogin  bool 		`json:"can_login" gorm:"default:false"`
	Password  string	`json:"password" gorm:"size:255"`
	Role	  string	`json:"role" gorm:"size:50;default:'employee'"`
	Salary	  float64	`json:"salary" gorm:"type:decimal(10,2);default:0"`
	DateJoined 	time.Time 		`json:"date_joined" `
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Station		Station			`json:"station" gorm:"foreignKey:StationID;references:ID;constraint:OnUpdate:CASCADE"`
	SalaryAdvance []SalaryAdvance `json:"salary_advance" gorm:"foreignKey:EmployeeID;references:ID;constraint:OnUpdate:CASCADE"`
	Payments []Payment `json:"payments" gorm:"foreignKey:EmployeeID;references:ID;constraint:OnUpdate:CASCADE"`
}


// SalaryAdvance represents an advance payment made to an employee
type SalaryAdvance struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	EmployeeID  uuid.UUID  `json:"employee_id" gorm:"type:varchar(36)"`
	Amount     	float64    `json:"amount" gorm:"type:decimal(10,2);not null"`
	Reason     	string     `json:"reason" gorm:"size:255"`
	DateRequested time.Time `json:"date_requested" gorm:"autoCreateTime"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Station struct{
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	Name      string      `json:"name"`
	Address   string      `json:"address"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	StationFuelProducts []StationFuelProduct `json:"station_fuel_prices" gorm:"foreignKey:StationID"`
	Employee	[]Employee	`json:"employees" gorm:"foreignKey:StationID;references:ID;constraint:OnUpdate:CASCADE"`
	Expenses    []Expenses 	   `json:"expenses" gorm:"foreignKey:StationID;references:ID;constraint:OnUpdate:CASCADE"`
	Tanks		[]Tank			`json:"tanks" gorm:"foreignKey:StationID;references:ID;constraint:OnUpdate:CASCADE"`
}

// Tank represents a tank at the station
type Tank struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	Name       string  `json:"name" gorm:"size:100"`
	Capacity   float64 `json:"capacity" gorm:"type:decimal(10,2);not null"`
	FuelProductID uuid.UUID `json:"fuel_product_id" gorm:"type:varchar(36);not null"`
	FuelProduct   FuelProduct   `json:"fuel_product" gorm:"foreignKey:FuelProductID"` // ✅ Added this line
	StationID    uuid.UUID `json:"station_id" gorm:"type:varchar(36);not null"`
	Station       Station   `json:"station" gorm:"foreignKey:StationID"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Pumps      []Pump  `gorm:"many2many:tank_pumps;" json:"pumps"`
	Dippings	[]Dippings `json:"dippings" gorm:"foreignKey:TankID;references:ID;constraint:OnUpdate:CASCADE"`
}
// Pump represents a pump at the station
type Pump struct {
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	Name       string   `json:"name" gorm:"size:100"`
	StationID  uuid.UUID `json:"station_id" gorm:"type:varchar(36);not null"`
	Tanks      []Tank   `gorm:"many2many:tank_pumps;" json:"tanks"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Nozzles    []Nozzle `json:"nozzles" gorm:"foreignKey:PumpID;references:ID;constraint:OnUpdate:CASCADE"`
	Readings 	[]PumpReadings `json:"readings" gorm:"foreignKey:PumpID;references:ID;constraint:OnUpdate:CASCADE"`
	Sales		[]Sales `json:"sales" gorm:"foreignKey:PumpID;references:ID;constraint:OnUpdate:CASCADE"`
}
// Nozzle represents a nozzle at the station
type Nozzle struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	Number        string    `json:"number" gorm:"size:100"`
	PumpID      uuid.UUID `json:"pump_id" gorm:"type:varchar(36);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// FuelProduct represents a fuel product available at the station
type FuelProduct struct {
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	Name      string      `json:"name"`
	Description string    `json:"description"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	StationFuelProducts []StationFuelProduct `json:"station_fuel_prices" gorm:"foreignKey:FuelProductID"`

	Tanks		[]Tank			`json:"tanks" gorm:"foreignKey:FuelProductID;references:ID;constraint:OnUpdate:CASCADE"`
}

type StationFuelProduct struct {
	ID         uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	StationID     uuid.UUID   `json:"station_id"`
	FuelProductID uuid.UUID   `json:"fuel_product_id"`
	UnitPrice     float64     `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	EffectiveFrom time.Time   `json:"effective_from" gorm:"not null"`
	CreatedAt     time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	Station     Station     `json:"station" gorm:"foreignKey:StationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	FuelProduct FuelProduct `json:"fuel_product" gorm:"foreignKey:FuelProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}



//fuel transaction 
type FuelTransaction struct {
    ID           uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
    FuelProductID uuid.UUID `json:"fuel_product_id" gorm:"type:varchar(36)"`
    StationID    uuid.UUID `json:"station_id" gorm:"type:varchar(36);not null"`
    Type         string    `json:"type" gorm:"type:enum('sale','supply','dipping','adjustment');not null"`
    Quantity     float64   `json:"quantity" gorm:"not null"`
    PreviousLevel float64  `json:"previous_level" gorm:"not null"`
    NewLevel     float64   `json:"new_level" gorm:"not null"`
    ReferenceID  uuid.UUID `json:"reference_id" gorm:"type:varchar(36)"` // Links to sale/supply record
    CreatedBy    uuid.UUID `json:"created_by" gorm:"type:varchar(36);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}


// Supplier represents a supplier of fuel products
type Supplier struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	Name        string    `json:"name" gorm:"size:100"`
	Description string    `json:"description" gorm:"size:255"`
	ContactName string    `json:"contact_name" gorm:"size:100"`
	PhoneNumber string    `json:"phone_number" gorm:"size:15"`
	Email       string    `json:"email" gorm:"size:100"`
	Address     string    `json:"address" gorm:"size:255"`
	CreditBalance float64 `json:"credit_balance" gorm:"type:decimal(10,2);default:0"` // Tracks overpayments
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Debts SupplierDebt `json:"supplier_debts" gorm:"foreignKey:SupplierID;references:ID;constraint:OnUpdate:CASCADE,"`
}

//supplies
type Supply struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	CarNumber     string    `json:"car_number" gorm:"size:20"`                          // Optional vehicle ID
	SupplierID    uuid.UUID `json:"supplier_id" gorm:"type:char(36);not null"`         // Foreign key to Supplier
	StationID     uuid.UUID `json:"station_id" gorm:"type:char(36);not null"`          // Foreign key to Station
	TankID          uuid.UUID  `json:"tank_id" gorm:"type:varchar(36);not null;"`
	EmployeeID    uuid.UUID `json:"employee_id" gorm:"type:char(36);not null"`         // Recorded by which employee
	ReferenceNo   string    `json:"reference_no" gorm:"size:50"`                       // Invoice or PO number
	FuelProductID uuid.UUID `json:"fuel_product_id" gorm:"type:char(36);not null"`     // FK to FuelProduct
	FuelProduct   FuelProduct   `json:"fuel_product" gorm:"foreignKey:FuelProductID;references:ID"`

	Quantity      float64   `json:"quantity" gorm:"type:decimal(10,2);not null"`       // Litres or gallons
	UnitPrice     float64   `json:"unit_price" gorm:"type:decimal(10,2);not null"`     // Cost per unit
	TotalAmount   float64   `json:"total_amount" gorm:"type:decimal(10,2);not null"`   // Quantity × UnitPrice
	DeliveryDate  time.Time `json:"delivery_date" gorm:"not null"`                     // Date delivered
	IsPaid        bool      `json:"is_paid" gorm:"default:false"`                      // Settled or not
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`                  // Auto timestamp
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}


// SupplierDebts represents debts owed to suppliers
type SupplierDebt struct {
	ID              uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	SupplierID      uuid.UUID `json:"supplier_id" gorm:"type:char(36);not null"`
	SupplyID        *uuid.UUID `json:"supply_id" gorm:"type:char(36)"` // Optional if transaction type is "payment"
	Supply          *Supply        `json:"supply,omitempty" gorm:"foreignKey:SupplyID;references:ID"`

	TransactionType string    `json:"transaction_type" gorm:"size:20;not null"` // "supply", "payment", etc.
	Amount          float64   `json:"amount" gorm:"type:decimal(10,2);not null"` 
	RunningBalance  float64   `json:"running_balance" gorm:"type:decimal(10,2);not null"` // After this transaction
	Notes           string    `json:"notes" gorm:"size:255"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

//supplier payments
type SupplierPayment struct {
    ID          uuid.UUID `json:"id" gorm:"primaryKey"`
    SupplierID  uuid.UUID `json:"supplier_id" gorm:""`
    Amount      float64   `json:"amount" gorm:"type:decimal(10,2)"`
    PaymentDate time.Time `json:"payment_date"`
    Method      string    `json:"method" gorm:"size:20"` // "cash", "transfer", etc.
    Reference   string    `json:"reference" gorm:"size:100"` // Payment reference
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Sales represents sales of fuel products
type Sales struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	EmployeeID  uuid.UUID `json:"employee_id" gorm:"type:varchar(36);"`
	PumpID     uuid.UUID `json:"pump_id" gorm:"type:varchar(36);not null"`
	LitersSold    float64       `json:"liters_sold" gorm:"type:decimal(10,2);not null"`
	PricePerLiter float64       `json:"price_per_liter" gorm:"type:decimal(10,2);not null"`
	TotalAmount   float64       `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Dippings represents the dippings of fuel products at a station
type Dippings struct {
	ID              uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	TankID          uuid.UUID  `json:"tank_id" gorm:"type:varchar(36);not null;"`
	DippingDate     time.Time `json:"dipping_date" gorm:"autoCreateTime"`
	OpeningDip      float64    `json:"opening_dip" gorm:"type:decimal(10,2);not null"`
	ClosingDip      float64    `json:"closing_dip" gorm:"type:decimal(10,2);not null"`
	OpeningMeter   float64    `json:"opening_meter" gorm:"type:decimal(10,2);not null"`
	ClosingMeter   float64    `json:"closing_meter" gorm:"type:decimal(10,2);not null"`
	LitersDispensed float64    `json:"liters_dispensed" gorm:"type:decimal(10,2);not null"`
	AmountSupplied float64    `json:"amount_supplied" gorm:"type:decimal(10,2);not null"`
	Deviation	  float64    `json:"deviation" gorm:"type:decimal(10,2);not null"` // Difference dip and sales
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Tank		Tank			`json:"tanks" gorm:"foreignKey:TankID;references:ID;constraint:OnUpdate:CASCADE"`
}

// Expenses represents expenses incurred by the station
type Expenses struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	StationID   uuid.UUID `json:"station_id" gorm:"type:varchar(36);"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Description string    `json:"description" gorm:"size:255"`
	ExpenseType string    `json:"expense_type" gorm:"size:100;not null"`
	ExpenseDate time.Time `json:"expense_date" gorm:"autoCreateTime"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
// Payments represents payments made to employees
type Payment struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36)"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Description string    `json:"description" gorm:"size:255"`
	PaymentDate time.Time `json:"payment_date" gorm:"autoCreateTime"`
	PaidMonth   string    `json:"paid_month" gorm:"size:50;not null"`
	Status		string    `json:"status" gorm:"size:50;not null default:'unpaid'"`
	EmployeeID  uuid.UUID `json:"employee_id" gorm:"type:varchar(36);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type PumpReadings struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);"`
	PumpID      uuid.UUID `json:"pump_id" gorm:"type:varchar(36);not null"`
	ReadingDate  time.Time `json:"reading_date" gorm:"autoCreateTime"`
	BusinessDay  time.Time `json:"business_day" gorm:"type:date;index"`
	Shift       string    `json:"shift" gorm:"size:50;not null"`
	OpeningMeter float64    `json:"opening_meter" gorm:"type:decimal(10,2);not null"`
	ClosingMeter float64    `json:"closing_meter" gorm:"type:decimal(10,2);not null"`
	LitersDispensed float64    `json:"liters_dispensed" gorm:"type:decimal(10,2);not null"`
	OpeningSalesAmount float64   `json:"opening_sales_amount" gorm:"type:decimal(10,2);not null"`
	ClosingSalesAmount float64   `json:"closing_sales_amount" gorm:"type:decimal(10,2);not null"`
	TotalSalesAmount  float64   `json:"total_sales_amount" gorm:"type:decimal(10,2);not null"`
	UnitPrice	 float64   `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	RecordedBy  uuid.UUID `json:"recorded_by" gorm:"type:varchar(36);not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Pump        Pump      `json:"pump" gorm:"foreignKey:PumpID;references:ID;constraint:OnUpdate:CASCADE"`
}

type FuelStock struct {
    ID            uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
    TankID        uuid.UUID `json:"tank_id" gorm:"not null;unique"`
    FuelProductID uuid.UUID `json:"fuel_product_id" gorm:"not null"`
    StationID     uuid.UUID `json:"station_id" gorm:"not null"`
    CurrentVolume float64   `json:"current_volume" gorm:"type:decimal(10,2)"`
    LastUpdated   time.Time `json:"last_updated" gorm:"autoUpdateTime"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`

    // Relationships
    FuelProduct   FuelProduct `json:"fuel_product" gorm:"foreignKey:FuelProductID;references:ID"`
    Tank          Tank        `json:"tank" gorm:"foreignKey:TankID;references:ID"`
    Station       Station     `json:"station" gorm:"foreignKey:StationID;references:ID"`
}

type Customer struct{
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	PhoneNumber string    `json:"phone_number" gorm:"size:15;not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
type CustomerCredit struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	CustomerID  uuid.UUID `json:"customer_id" gorm:"type:varchar(36);not null"`
	Date		time.Time `json:"date" gorm:"autoCreateTime"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
//customer credits payments
type CustomerCreditPayment struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	CustomerID  uuid.UUID `json:"customer_id" gorm:"type:varchar(36);not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentDate time.Time `json:"payment_date" gorm:"autoCreateTime"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`	
}

//accounts
type DailyAccounts struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	StationID   uuid.UUID `json:"station_id" gorm:"type:varchar(36);not null"`
	Station		Station   `json:"station" gorm:"foreignKey:StationID;references:ID;constraint:OnUpdate:CASCADE"`
	TotalSalesAmount float64 `json:"total_sales_amount" gorm:"type:decimal(10,2);not null"`
	BusinessDay		time.Time `json:"business_day" gorm:"type:date;index"` // Date for which accounts are being recorded
	TotalExpenses float64 `json:"total_expenses" gorm:"type:decimal(10,2);not null"`
	DebtTaken float64 `json:"debt_taken" gorm:"type:decimal(10,2);not null"`
	DebtPaid float64 `json:"debt_paid" gorm:"type:decimal(10,2);not null"`
	Mpesa float64 `json:"mpesa" gorm:"type:decimal(10,2);not null"`
	Bank float64 `json:"bank" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt	*gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

//before save supply hook
func (s *Supply) BeforeSave(tx *gorm.DB) (err error) {
	s.TotalAmount = s.Quantity * s.UnitPrice
	return nil
}

