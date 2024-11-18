package handlers

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type BusyOrder struct {
	ServiceID string    `gorm:"column:service_id" json:"service_id"`
	OrderDate time.Time `gorm:"column:order_date" json:"order_date"`
}

func ValidateOrder(db *gorm.DB, order ClientOrder, receiptDate time.Time) (map[string]interface{}, error) {
	if order.OrderDate.Before(time.Now()) {
		return nil, fmt.Errorf("order date can't be in the past")
	}

	fmt.Printf("Order date: %v, Receipt date: %v\n", order.OrderDate, receiptDate)

	workerID, err := SetWorker(db, order.OrderDate, receiptDate)
	if err != nil {
		return nil, fmt.Errorf("error assigning worker: %v", err)
	}

	return map[string]interface{}{
		"worker_id": workerID,
	}, nil
}
