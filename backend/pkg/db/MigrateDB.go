package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
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
		&BookingToOrder{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

func InsertUniqueDB(db *gorm.DB) error {
	if err := insertServices(db); err != nil {
		return err
	}

	if err := insertMaterials(db); err != nil {
		return err
	}

	if err := insertEquipments(db); err != nil {
		return err
	}

	if err := insertServiceRequirements(db); err != nil {
		return err
	}

	if err := insertClients(db); err != nil {
		return err
	}

	if err := insertEmployees(db); err != nil {
		return err
	}

	if err := insertOrders(db); err != nil {
		return err
	}

	if err := insertBookings(db); err != nil {
		return err
	}

	if err := insertMaterialExpenditures(db); err != nil {
		return err
	}

	if err := insertMaterialPurchases(db); err != nil {
		return err
	}

	if err := insertOrderContents(db); err != nil {
		return err
	}

	if err := insertBookingToOrders(db); err != nil {
		return err
	}

	return nil
}

func insertServices(db *gorm.DB) error {
	for _, service := range Services {
		if service.Name == "" {
			log.Println("Пропуск пустого названия услуги")
			continue
		}

		if err := db.FirstOrCreate(&service, Service{Name: service.Name}).Error; err != nil {
			return fmtError("услуги", err)
		}
	}
	return nil
}

func insertMaterials(db *gorm.DB) error {
	for _, material := range Materials {
		if material.MaterialName == "" {
			log.Println("Пропуск пустого названия материала")
			continue
		}

		if err := db.FirstOrCreate(&material, Material{MaterialName: material.MaterialName}).Error; err != nil {
			return fmtError("материала", err)
		}
	}
	return nil
}

func insertEquipments(db *gorm.DB) error {
	for _, equipment := range Equipments {
		if equipment.Brand == "" || equipment.Model == "" {
			log.Println("Пропуск оборудования с пустым брендом или моделью")
			continue
		}

		if err := db.FirstOrCreate(&equipment, Equipment{Brand: equipment.Brand, Model: equipment.Model}).Error; err != nil {
			return fmtError("оборудования", err)
		}
	}
	return nil
}

func insertServiceRequirements(db *gorm.DB) error {
	for _, sr := range ServiceRequirements {

		if sr.ServiceID == 0 || sr.EquipmentID == 0 {
			log.Println("Пропуск ServiceRequirement с пустыми ServiceID или EquipmentID")
			continue
		}

		if err := db.FirstOrCreate(&sr, ServiceRequirement{ServiceID: sr.ServiceID, EquipmentID: sr.EquipmentID}).Error; err != nil {
			return fmtError("ServiceRequirement", err)
		}
	}
	return nil
}

func insertClients(db *gorm.DB) error {
	for _, client := range Clients {
		if client.FullName == "" || client.PhoneNumber == "" || client.Email == "" {
			log.Println("Пропуск клиента с пустыми полями")
			continue
		}

		if err := db.FirstOrCreate(&client, Client{PhoneNumber: client.PhoneNumber, Email: client.Email}).Error; err != nil {
			return fmtError("клиента", err)
		}
	}
	return nil
}

func insertEmployees(db *gorm.DB) error {
	for _, employee := range Employees {
		if employee.FullName == "" || employee.Position == "" || employee.PassportData == "" || employee.PhoneNumber == "" {
			log.Println("Пропуск сотрудника с пустыми полями")
			continue
		}

		if err := db.FirstOrCreate(&employee, Employee{PassportData: employee.PassportData, PhoneNumber: employee.PhoneNumber}).Error; err != nil {
			return fmtError("сотрудника", err)
		}
	}
	return nil
}

func insertOrders(db *gorm.DB) error {
	for _, order := range Orders {
		if order.ServiceName == "" {
			log.Println("Пропуск заказа с пустым названием услуги")
			continue
		}

		if err := db.FirstOrCreate(&order, Order{OrderID: order.OrderID}).Error; err != nil {
			return fmtError("заказа", err)
		}
	}
	return nil
}

func insertBookings(db *gorm.DB) error {
	for _, booking := range Bookings {
		if booking.BookerFullName == "" || booking.BookingType == "" {
			log.Println("Пропуск бронирования с пустыми полями")
			continue
		}

		if err := db.FirstOrCreate(&booking, Booking{BookingID: booking.BookingID}).Error; err != nil {
			return fmtError("бронирования", err)
		}
	}
	return nil
}

func insertMaterialExpenditures(db *gorm.DB) error {
	for _, me := range MaterialExpenditures {
		if me.MaterialID == 0 || me.Quantity == 0 {
			log.Println("Пропуск MaterialExpenditure с пустыми полями")
			continue
		}

		if err := db.FirstOrCreate(&me, MaterialExpenditure{ExpenditureID: me.ExpenditureID}).Error; err != nil {
			return fmtError("MaterialExpenditure", err)
		}
	}
	return nil
}

func insertMaterialPurchases(db *gorm.DB) error {
	for _, mp := range MaterialPurchases {
		if mp.MaterialID == 0 || mp.Supplier == "" {
			log.Println("Пропуск MaterialPurchase с пустыми полями")
			continue
		}

		if err := db.FirstOrCreate(&mp, MaterialPurchase{PurchaseID: mp.PurchaseID}).Error; err != nil {
			return fmtError("MaterialPurchase", err)
		}
	}
	return nil
}

func insertOrderContents(db *gorm.DB) error {
	for _, oc := range OrderContents {
		if oc.OrderID == 0 || oc.ServiceID == 0 {
			log.Println("Пропуск OrderContent с пустыми полями")
			continue
		}

		if err := db.FirstOrCreate(&oc, OrderContent{OrderID: oc.OrderID, ServiceID: oc.ServiceID}).Error; err != nil {
			return fmtError("OrderContent", err)
		}
	}
	return nil
}

func insertBookingToOrders(db *gorm.DB) error {
	for _, bto := range BookingToOrders {
		if bto.BookingID == 0 || bto.EmployeeID == 0 || bto.ClientID == 0 {
			log.Println("Пропуск BookingToOrder с пустыми полями")
			continue
		}

		if err := db.FirstOrCreate(&bto, BookingToOrder{ID: bto.ID}).Error; err != nil {
			return fmtError("BookingToOrder", err)
		}
	}
	return nil
}

func fmtError(model string, err error) error {
	return fmt.Errorf("ошибка при создании %s: %w", model, err)
}
