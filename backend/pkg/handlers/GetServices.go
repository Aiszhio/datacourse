package handlers

import (
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetServices(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var services []db2.Service

		if err := db.Table("services").Select("service_id, price, name").Find(&services).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error while fetching services",
			})
		}
		return c.Status(fiber.StatusOK).JSON(services)
	}
}
