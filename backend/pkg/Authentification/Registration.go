package Authentification

import (
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/Aiszhio/datacourse.git/pkg/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Registration(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var client db2.Client

		if err := c.BodyParser(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга данных",
			})
		}

		if err := handlers.ValidateClient(client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Table("clients").Create(&client).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Вы успешно зарегистрировались!",
		})
	}
}
