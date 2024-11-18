package handlers

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

func SetWorker(db *gorm.DB, orderStart time.Time, orderEnd time.Time) (int, error) {
	allowedWorkerIDs := []int{1, 2, 4, 6, 8}

	var busyWorkerIDs []int
	err := db.Table("orders").
		Select("DISTINCT employee_id").
		Where("order_date < ? AND receipt_date > ?", orderEnd, orderStart).
		Where("employee_id IN ?", allowedWorkerIDs).
		Pluck("employee_id", &busyWorkerIDs).Error
	if err != nil {
		return 0, fmt.Errorf("ошибка при поиске занятых работников: %v", err)
	}

	availableWorkerIDs := difference(allowedWorkerIDs, busyWorkerIDs)

	if len(availableWorkerIDs) == 0 {
		return 0, fmt.Errorf("нет доступных работников на указанный временной промежуток")
	}

	availableWorkerID := availableWorkerIDs[0]

	fmt.Printf("Доступный ID работника: %v\n", availableWorkerID)
	return availableWorkerID, nil
}

func difference(all []int, busy []int) []int {
	busyMap := make(map[int]bool)
	for _, id := range busy {
		busyMap[id] = true
	}

	var diff []int
	for _, id := range all {
		if !busyMap[id] {
			diff = append(diff, id)
		}
	}
	return diff
}
