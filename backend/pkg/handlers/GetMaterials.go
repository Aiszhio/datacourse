package handlers

import (
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetMaterials(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var materialsList []database.Material

		err := db.Table("materials").Find(&materialsList).Error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"materials": materialsList,
		})
	}
}
