package handlers

import (
	"errors"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"regexp"
	"strconv"
)

type EquipmentToUpdate struct {
	Type  string `gorm:"column:type;not null" json:"type"`
	Brand string `gorm:"column:brand;not null" json:"brand"`
	Model string `gorm:"column:model;not null" json:"model"`
}

func UpdateEquipment(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		equipmentIdStr := c.Params("equipment_id")

		if equipmentIdStr == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		equipmentId, err := strconv.Atoi(equipmentIdStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err,
			})
		}

		var updatedEquipment EquipmentToUpdate
		var existingEquipment db2.Equipment

		if err = db.Table("equipment").Where("equipment_id = ?", equipmentId).First(&existingEquipment).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Такого оборудования не существует",
			})
		}

		if err = c.BodyParser(&updatedEquipment); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга данных",
			})
		}

		if err = ValidateEquipment(updatedEquipment); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err = db.Model(&existingEquipment).Updates(updatedEquipment).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка сохранения данных в БД",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Данные об оборудовании были успешно обновлены!",
			"service": updatedEquipment,
		})
	}
}

func ValidateEquipment(updatedEquipment EquipmentToUpdate) error {

	modelRegex := "^[a-zA-Z0-9- ]+$"
	if !regexp.MustCompile(modelRegex).MatchString(updatedEquipment.Model) {
		return errors.New("Модель должна содержать только буквы, цифры, пробелы и дефисы.")
	}

	if len(updatedEquipment.Model) < 3 {
		return errors.New("Модель должна содержать хотя бы 3 символа.")
	}

	return nil
}
