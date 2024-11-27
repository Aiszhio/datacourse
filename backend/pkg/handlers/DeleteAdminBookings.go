package handlers

import (
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteAdminBooking(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bookingID := c.Params("bookingID")

		tx := db.Begin()

		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error,
			})
		}

		if err := tx.Where("booking_id = ?", bookingID).Delete(&db2.Booking{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		if err := tx.Where("booking_id = ?", bookingID).Delete(&db2.BookingToOrder{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		if err := tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Бронирование успешно удалено",
		})
	}
}
