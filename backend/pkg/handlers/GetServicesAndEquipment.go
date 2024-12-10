package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Equipment struct {
	ID    int    `gorm:"column:equipment_id;primaryKey;autoIncrement" json:"equipment_id"`
	Type  string `gorm:"column:type" json:"type"`
	Brand string `gorm:"column:brand;not null" json:"brand"`
	Model string `gorm:"column:model;not null" json:"model"`
}

type Service struct {
	ID                int         `gorm:"column:service_id;primaryKey;autoIncrement" json:"service_id"`
	Price             float64     `gorm:"column:price;not null" json:"price"`
	Name              string      `gorm:"column:name" json:"name"`
	RequiredEquipment []Equipment `gorm:"many2many:service_requirements;"`
}

func GetServicesAndEquipment(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var services []Service

		err := db.Preload("RequiredEquipment").Find(&services).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при получении данных услуг и оборудования",
			})
		}

		var equipment []Equipment
		err = db.Table("equipment").Find(&equipment).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		fmt.Println(services)

		return c.JSON(fiber.Map{
			"services":  services,
			"equipment": equipment,
		})
	}
}
