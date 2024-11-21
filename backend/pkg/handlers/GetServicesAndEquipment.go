package handlers

import (
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetServicesAndEquipment(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var serviceList []database.Service
		var equipmentList []database.Equipment

		err := db.Table("services").Find(&serviceList).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		err = db.Table("equipment").Find(&equipmentList).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.JSON(fiber.Map{
			"equipment": equipmentList,
			"services":  serviceList,
		})
	}
}
