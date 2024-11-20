package handlers

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func ValidateService(id int, db *gorm.DB) error {
	var equipmentID int

	if err := db.Table("service_requirements").Select("equipment_id").
		Where("id = ?", id).Scan(&equipmentID).Error; err != nil {
		fmt.Println("Error in ValidateService: ", err)
		return err
	}

	if equipmentID == 0 {
		return errors.New("No equipment for service")
	}

	return nil
}
