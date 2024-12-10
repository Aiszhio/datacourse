package db

import (
	"time"
)

type Client struct {
	ClientID    int    `gorm:"column:client_id;primaryKey;autoIncrement"`
	FullName    string `gorm:"column:full_name;not null" json:"FullName"`
	PhoneNumber string `gorm:"column:phone_number;not null;unique" json:"PhoneNumber"`
	Email       string `gorm:"column:email;not null;unique" json:"Email"`
}

type Order struct {
	OrderID     int       `gorm:"column:order_id;primaryKey" json:"id"`
	ClientID    int       `gorm:"column:client_id" json:"clientId"`
	EmployeeID  int       `gorm:"column:employee_id" json:"employeeId"`
	ServiceName string    `gorm:"column:service_name" json:"service"`
	OrderDate   time.Time `gorm:"column:order_date" json:"orderDate"`
	ReceiptDate time.Time `gorm:"column:receipt_date" json:"receiveDate"`
}

type Employee struct {
	EmployeeID   int       `gorm:"column:employee_id;primaryKey;autoIncrement" json:"employee_id"`
	FullName     string    `gorm:"column:full_name;not null" json:"full_name"`
	Position     string    `gorm:"column:position;not null" json:"position"`
	HireDate     time.Time `gorm:"column:hire_date;not null" json:"hire_date"`
	BirthDate    time.Time `gorm:"column:birth_date;not null" json:"birth_date"`
	PassportData string    `gorm:"column:passport_data;not null;unique" json:"passport_data"`
	PhoneNumber  string    `gorm:"column:phone_number;not null;unique" json:"phone_number"`
	Status       string    `gorm:"column:status;not null;default:'Работает'" json:"status"`
}

type Service struct {
	ServiceID int     `gorm:"column:service_id;primaryKey;autoIncrement" json:"service_id"`
	Price     float64 `gorm:"column:price;not null" json:"price"`
	Name      string  `gorm:"column:name" json:"name"`
}

type Booking struct {
	BookingID      int       `gorm:"column:booking_id;primaryKey;autoIncrement" json:"id"`
	BookingType    string    `gorm:"column:booking_type;not null" json:"booking_type"`
	OrderID        int       `gorm:"column:order_id;not null;autoIncrement" json:"order_id"`
	BookingTime    time.Time `gorm:"column:booking_time;not null" json:"booking_time"`
	BookerFullName string    `gorm:"column:booker_full_name;not null" json:"booker_full_name"`
	ClientID       int       `gorm:"column:client_id;" json:"client_id"`
}

type Material struct {
	MaterialID   int    `gorm:"column:material_id;primaryKey;autoIncrement" json:"material_id"`
	MaterialName string `gorm:"column:material_name;not null;unique" json:"material_name"`
	Quantity     int    `gorm:"column:quantity;default:0" json:"quantity"`
}

type MaterialExpenditure struct {
	ExpenditureID   int       `gorm:"column:expenditure_id;primaryKey;autoIncrement" json:"expenditure_id"`
	MaterialID      int       `gorm:"column:material_id;not null" json:"material_id"`
	ExpenditureDate time.Time `gorm:"column:expenditure_date;not null" json:"expenditure_date"`
	Quantity        int       `gorm:"column:quantity;not null" json:"quantity"`
}

type MaterialPurchase struct {
	PurchaseID int       `gorm:"column:purchase_id;primaryKey;autoIncrement" json:"purchase_id"`
	MaterialID int       `gorm:"column:material_id;not null" json:"material_id"`
	Cost       float64   `gorm:"column:cost;not null" json:"cost"`
	Supplier   string    `gorm:"column:supplier;not null" json:"supplier"`
	Quantity   int       `gorm:"column:quantity;not null" json:"quantity"`
	SupplyDate time.Time `gorm:"column:supply_date;not null" json:"supply_date"`
}

type Equipment struct {
	EquipmentID int    `gorm:"column:equipment_id;primaryKey;autoIncrement" json:"equipment_id"`
	Type        string `gorm:"column:type" json:"type"`
	Brand       string `gorm:"column:brand;not null" json:"brand"`
	Model       string `gorm:"column:model;not null" json:"model"`
}

type ServiceRequirement struct {
	ServiceID   int `gorm:"column:service_id;primaryKey"`
	EquipmentID int `gorm:"column:equipment_id;primaryKey"`
}

type OrderContent struct {
	OrderID   int `gorm:"column:order_id;primaryKey"`
	ServiceID int `gorm:"column:service_id;primaryKey"`
}

type BookingToOrder struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BookingID  int       `gorm:"column:booking_id;not null" json:"booking_id"`
	EmployeeID int       `gorm:"column:employee_id;not null" json:"employee_id"`
	ClientID   int       `gorm:"column:client_id;not null" json:"client_id"`
	OrderDate  time.Time `gorm:"column:order_date;not null" json:"orderDate"`
	OrderEnd   time.Time `gorm:"column:order_end;not null" json:"orderEnd"`
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

var BookingToOrders = []BookingToOrder{}
