package db

import (
	"time"
)

type Client struct {
	ClientID    int    `gorm:"column:client_id;primaryKey;autoIncrement"`
	FullName    string `gorm:"column:full_name;not null"`
	PhoneNumber string `gorm:"column:phone_number;not null;unique"`
	Email       string `gorm:"column:email;not null;unique"`
}

type Order struct {
	OrderID     int       `gorm:"column:order_id;primaryKey;autoIncrement"`
	ClientID    int       `gorm:"column:client_id;not null"`
	EmployeeID  int       `gorm:"column:employee_id;not null"`
	ServiceName string    `gorm:"column:service_name;not null"`
	OrderDate   time.Time `gorm:"column:order_date;not null"`
	ReceiptDate time.Time `gorm:"column:receipt_date;not null"`
}

type Employee struct {
	EmployeeID   int       `gorm:"column:employee_id;primaryKey;autoIncrement"`
	FullName     string    `gorm:"column:full_name;not null"`
	Position     string    `gorm:"column:position;not null"`
	HireDate     time.Time `gorm:"column:hire_date;not null"`
	BirthDate    time.Time `gorm:"column:birth_date;not null"`
	PassportData string    `gorm:"column:passport_data;not null;unique"`
	PhoneNumber  string    `gorm:"column:phone_number;not null;unique"`
}

type Service struct {
	ServiceID int     `gorm:"column:service_id;primaryKey;autoIncrement"`
	Price     float64 `gorm:"column:price;not null"`
	Name      string  `gorm:"column:name"`
}

type Booking struct {
	BookingID      int       `gorm:"column:booking_id;primaryKey;autoIncrement"`
	BookingType    string    `gorm:"column:booking_type;not null"`
	OrderID        int       `gorm:"column:order_id;not null"`
	BookingTime    time.Time `gorm:"column:booking_time;not null"`
	BookerFullName string    `gorm:"column:booker_full_name;not null"`
}

type Material struct {
	MaterialID   int    `gorm:"column:material_id;primaryKey;autoIncrement"`
	MaterialName string `gorm:"column:material_name;not null;unique"`
}

type MaterialExpenditure struct {
	ExpenditureID   int       `gorm:"column:expenditure_id;primaryKey;autoIncrement"`
	MaterialID      int       `gorm:"column:material_id;not null"`
	ExpenditureDate time.Time `gorm:"column:expenditure_date;not null"`
	Quantity        int       `gorm:"column:quantity;not null"`
}

type MaterialPurchase struct {
	PurchaseID int       `gorm:"column:purchase_id;primaryKey;autoIncrement"`
	MaterialID int       `gorm:"column:material_id;not null"`
	Cost       float64   `gorm:"column:cost;not null"`
	Supplier   string    `gorm:"column:supplier;not null"`
	Quantity   int       `gorm:"column:quantity;not null"`
	SupplyDate time.Time `gorm:"column:supply_date;not null"`
}

type Equipment struct {
	EquipmentID int    `gorm:"column:equipment_id;primaryKey;autoIncrement"`
	Brand       string `gorm:"column:brand;not null"`
	Model       string `gorm:"column:model;not null"`
}

type ServiceRequirement struct {
	ServiceID   int `gorm:"column:service_id;primaryKey;autoIncrement"`
	EquipmentID int `gorm:"column:equipment_id;primaryKey"`
}

type OrderContent struct {
	OrderID   int `gorm:"column:order_id;primaryKey;autoIncrement"`
	ServiceID int `gorm:"column:service_id;primaryKey"`
}

var Services = []Service{}

var Materials = []Material{}

var Equipments = []Equipment{}

var ServiceRequirements = []ServiceRequirement{}

var Clients = []Client{}

var Orders = []Order{}

var Employees = []Employee{}

var Bookings = []Booking{}

var MaterialExpenditures = []MaterialExpenditure{}

var MaterialPurchases = []MaterialPurchase{}

var OrderContents = []OrderContent{}
