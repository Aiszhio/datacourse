package handlers

import (
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetAdminBookings(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var bookingsList []database.Booking

		err := db.Table("bookings").Find(&bookingsList).Error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"bookings": bookingsList,
		})
	}
}
