package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"gorm.io/gorm"
	"time"
)

func ValidateOrderEquipment(order db2.Order, db *gorm.DB) error {
	var requiredEquipment []int

	var service db2.Service
	if err := db.Table("services").Where("name = ?", order.ServiceName).First(&service).Error; err != nil {
		return fmt.Errorf("Ошибка при поиске услуги")
	}

	if err := db.Table("service_requirements").Where("service_id = ?", service.ServiceID).Pluck("equipment_id", &requiredEquipment).Error; err != nil {
		return fmt.Errorf("Ошибка при получении требуемого оборудования")
	}

	if len(requiredEquipment) == 0 {
		return fmt.Errorf("Для выбранной услуги не требуется оборудование")
	}

	var overlappingOrderIDs []int
	if err := db.Table("orders").
		Where("orders.receipt_date > ? AND orders.order_date < ?", order.OrderDate, order.ReceiptDate).
		Pluck("order_id", &overlappingOrderIDs).Error; err != nil {
		return fmt.Errorf("Ошибка при поиске перекрывающихся заказов")
	}

	if len(overlappingOrderIDs) == 0 {
		return nil
	}

	var overlappingServiceIDs []int
	if err := db.Table("order_contents").
		Where("order_contents.order_id IN ?", overlappingOrderIDs).
		Pluck("service_id", &overlappingServiceIDs).Error; err != nil {
		return fmt.Errorf("Ошибка при получении услуги")
	}

	if len(overlappingServiceIDs) == 0 {
		return nil
	}

	var occupiedEquipment []int
	if err := db.Table("service_requirements").
		Where("service_id IN ?", overlappingServiceIDs).
		Pluck("equipment_id", &occupiedEquipment).Error; err != nil {
		return fmt.Errorf("Ошибка при получении занятых оборудований")
	}

	occupiedSet := make(map[int]bool)
	for _, eq := range occupiedEquipment {
		occupiedSet[eq] = true
	}

	var unavailable []int
	for _, eq := range requiredEquipment {
		if occupiedSet[eq] {
			unavailable = append(unavailable, eq)
		}
	}

	if len(unavailable) > 0 {
		return fmt.Errorf("Нет свободного оборудования")
	}

	return nil
}

func ValidateDeleteEquipment(equipmentID int, db *gorm.DB) error {
	now := time.Now()

	var serviceIDs []int
	if err := db.Table("service_requirements").
		Where("equipment_id = ?", equipmentID).
		Pluck("service_id", &serviceIDs).Error; err != nil {
		return fmt.Errorf("Ошибка при получении идентификаторов оборудования")
	}

	if len(serviceIDs) == 0 {
		return nil
	}

	var orderIDs []int
	if err := db.Table("order_contents").
		Where("service_id IN ?", serviceIDs).
		Pluck("order_id", &orderIDs).Error; err != nil {
		return fmt.Errorf("Ошибка при получении идентификаторов для сервисов")
	}

	if len(orderIDs) == 0 {
		return nil
	}

	var count int64
	if err := db.Table("orders").
		Where("order_id IN ? AND receipt_date > ?", orderIDs, now).
		Count(&count).Error; err != nil {
		return fmt.Errorf("Ошибка при проверке занятости оборудования в заказах")
	}

	if count > 0 {
		return fmt.Errorf("Оборудование используется в активных или будущих заказах")
	}

	return nil
}
