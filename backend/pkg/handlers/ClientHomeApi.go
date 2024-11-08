package handlers

import (
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetUserData(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var orders []database.Order
		for _, order := range orders {
			fmt.Println("Order:", order)
		}
		return nil
	}
}
