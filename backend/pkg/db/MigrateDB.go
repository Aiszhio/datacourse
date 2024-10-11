package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MigrateDB() error {
	dbLink := "postgres://postgres:1234@postgres:5432/parlourdb?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dbLink), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = db.AutoMigrate(
		&Client{},
		&Order{},
		&Employee{},
		&Service{},
		&Booking{},
		&Material{},
		&MaterialExpenditure{},
		&MaterialPurchase{},
		&Equipment{},
		&ServiceRequirement{},
		&OrderContent{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

func InsertUniqueDB(db *gorm.DB) error {
	for _, service := range Services {
		var existingService Service
		if service.Name == "" {
			fmt.Println("Skipping empty service name!!!")
			continue
		}
		if err := db.FirstOrCreate(&existingService, Service{ServiceID: service.ServiceID}).Error; err != nil {
			return err
		}
	}

	for _, material := range Materials {
		var existingMaterial Material
		if material.MaterialName == "" {
			fmt.Println("Skipping empty material name")
			continue
		}
		if err := db.FirstOrCreate(&existingMaterial, Material{MaterialID: material.MaterialID}).Error; err != nil {
			return err
		}
	}

	for _, equipment := range Equipments {
		var existingEquipment Equipment
		if err := db.FirstOrCreate(&existingEquipment, Equipment{EquipmentID: equipment.EquipmentID}).Error; err != nil {
			return err
		}
	}

	for _, serviceRequirement := range ServiceRequirements {
		var existingServiceRequirement ServiceRequirement
		if err := db.FirstOrCreate(&existingServiceRequirement, ServiceRequirement{ServiceID: serviceRequirement.ServiceID}).Error; err != nil {
			return err
		}
	}

	for _, client := range Clients {
		var existingClient Client
		if err := db.FirstOrCreate(&existingClient, Client{ClientID: client.ClientID}).Error; err != nil {
			return err
		}
	}

	for _, order := range Orders {
		var existingOrder Order
		if err := db.FirstOrCreate(&existingOrder, Order{OrderID: order.OrderID}).Error; err != nil {
			return err
		}
	}

	for _, employee := range Employees {
		var existingEmployee Employee
		if err := db.FirstOrCreate(&existingEmployee, Employee{EmployeeID: employee.EmployeeID}).Error; err != nil {
			return err
		}
	}

	for _, booking := range Bookings {
		var existingBooking Booking
		if err := db.FirstOrCreate(&existingBooking, Booking{BookingID: booking.BookingID}).Error; err != nil {
			return err
		}
	}

	for _, materialExpenditure := range MaterialExpenditures {
		var existingMaterialExpenditure MaterialExpenditure
		if err := db.FirstOrCreate(&existingMaterialExpenditure, MaterialExpenditure{ExpenditureID: materialExpenditure.ExpenditureID}).Error; err != nil {
			return err
		}
	}

	for _, materialPurchase := range MaterialPurchases {
		var existingMaterialPurchase MaterialPurchase
		if err := db.FirstOrCreate(&existingMaterialPurchase, MaterialPurchase{PurchaseID: materialPurchase.PurchaseID}).Error; err != nil {
			return err
		}
	}

	for _, orderContent := range OrderContents {
		var existingOrderContent OrderContent
		if err := db.FirstOrCreate(&existingOrderContent, OrderContent{OrderID: orderContent.OrderID}).Error; err != nil {
			return err
		}
	}

	return nil
}
