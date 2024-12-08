package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AddClient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var clientList []db2.Client

		err := db.Table("clients").Find(&clientList).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Не нашли данные в таблице клиентов",
			})
		}

		if len(clientList) == 0 {
			fmt.Println("Ничего нет")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"clients": clientList,
		})
	}
}
