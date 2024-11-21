package handlers

import (
	"gorm.io/gorm"
)

func GetClientName(id int, db *gorm.DB) (string, error) {
	var clientName string

	err := db.Table("clients").Select("full_name").Where("client_id = ?", id).Scan(&clientName).Error
	if err != nil {
		return "", err
	}

	return clientName, nil
}

func GetWorkerName(id int, db *gorm.DB) (string, error) {
	var workerName string

	err := db.Table("employees").Select("full_name").Where("employee_id = ?", id).Scan(&workerName).Error
	if err != nil {
		return "", err
	}

	return workerName, nil
}
