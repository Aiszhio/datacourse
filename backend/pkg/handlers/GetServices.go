package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Service struct {
	ServiceID   int    `gorm:"column:service_id;primaryKey" json:"service_id"`
	Price       int    `gorm:"column:price;not null" json:"price"`
	ServiceName string `gorm:"column:name;not null" json:"service_name"`
}

func GetServices(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var services []Service

		if err := db.Table("services").Select("service_id, price ,name").Find(&services).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error while fetching services",
			})
		}
		return c.Status(fiber.StatusOK).JSON(services)
	}
}
