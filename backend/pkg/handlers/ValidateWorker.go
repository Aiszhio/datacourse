package handlers

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func SetWorker(db *gorm.DB, orderStart time.Time, orderEnd time.Time) (int, error) {
	var busyWorkerIDs []int
	err := db.Table("orders").
		Select("DISTINCT orders.employee_id").
		Joins("LEFT JOIN employees ON employees.employee_id = orders.employee_id").
		Where("orders.order_date < ? AND orders.receipt_date > ?", orderEnd, orderStart).
		Where("employees.position = ?", "Фотограф").
		Pluck("orders.employee_id", &busyWorkerIDs).Error
	if err != nil {
		return 0, fmt.Errorf("ошибка при поиске занятых работников: %v", err)
	}

	var allPhotographers []int
	err = db.Table("employees").
		Select("employees.employee_id").
		Where("employees.position = ?", "Фотограф").
		Pluck("employees.employee_id", &allPhotographers).Error
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении списка фотографов: %v", err)
	}

	freeWorkers, err := difference(busyWorkerIDs, allPhotographers)
	if err != nil {
		fmt.Println("Error difference busy workers")
		return 0, err
	}

	if len(freeWorkers) == 0 {
		fmt.Println("No workers available")
		return 0, nil
	}

	return freeWorkers[0], nil
}

func difference(busyWorkers, allWorkers []int) ([]int, error) {
	busyMap := make(map[int]struct{}, len(busyWorkers))
	for _, id := range busyWorkers {
		busyMap[id] = struct{}{}
	}

	var freeWorkers []int
	for _, id := range allWorkers {
		if _, found := busyMap[id]; !found {
			freeWorkers = append(freeWorkers, id)
		}
	}
	return freeWorkers, nil
}
