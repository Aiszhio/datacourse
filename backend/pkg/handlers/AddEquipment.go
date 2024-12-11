package handlers

import (
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddEquipment(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var equipment db2.Equipment

		if err := c.BodyParser(&equipment); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга данных",
			})
		}

		if err := db.Table("equipment").Create(&equipment).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка при добавлении оборудования в БД",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Оборудование успешно добавлено",
		})
	}
}
