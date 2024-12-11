package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

func ValidateService(input ServiceInput, db *gorm.DB) error {
	validNames := []string{"Печать", "Фотокнига", "Фотоальбом", "Фотосессия", "Съемка", "Аренда"}

	matched := false
	for _, name := range validNames {
		if strings.Contains(strings.ToLower(input.Name), strings.ToLower(name)) {
			matched = true
			break
		}
	}
	if !matched {
		return fmt.Errorf("Название услуги должно содержать одно из следующих слов: Печать, Фотокнига, Фотоальбом, Фотосессия, Съемка, Аренда")
	}

	priceStr := strconv.FormatFloat(input.Price, 'f', -1, 64)
	if len(priceStr) >= 7 {
		return fmt.Errorf("Вы ввели слишком большую стоимость")
	}

	if len(input.RequiredEquipment) > 0 {
		var count int64
		if err := db.Model(&Equipment{}).Where("equipment_id IN ?", input.RequiredEquipment).Count(&count).Error; err != nil {
			fmt.Println("Error in ValidateService while counting equipments: ", err)
			return fmt.Errorf("Ошибка при проверке оборудования")
		}

		if int(count) != len(input.RequiredEquipment) {
			return fmt.Errorf("Система не обнаружила часть оборудования")
		}
	}

	return nil
}
