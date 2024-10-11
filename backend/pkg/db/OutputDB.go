package db

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAllData(db *gorm.DB) (fiber.Map, error) {
	var clients []Client
	var orders []Order
	var employees []Employee
	var services []Service
	var bookings []Booking
	var materials []Material
	var materialExpenditures []MaterialExpenditure
	var materialPurchases []MaterialPurchase
	var equipments []Equipment
	var serviceRequirements []ServiceRequirement
	var orderContents []OrderContent

	if err := db.Find(&clients).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&orders).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&employees).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&services).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&bookings).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&materials).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&materialExpenditures).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&materialPurchases).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&equipments).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&serviceRequirements).Error; err != nil {
		return nil, err
	}
	if err := db.Find(&orderContents).Error; err != nil {
		return nil, err
	}

	return fiber.Map{
		"clients":               clients,
		"orders":                orders,
		"employees":             employees,
		"services":              services,
		"bookings":              bookings,
		"materials":             materials,
		"material_expenditures": materialExpenditures,
		"material_purchases":    materialPurchases,
		"equipments":            equipments,
		"service_requirements":  serviceRequirements,
		"order_contents":        orderContents,
	}, nil
}
