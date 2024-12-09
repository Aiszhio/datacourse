package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

func GetBookings(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var bookingsList []database.Booking

		userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Не удалось получить данные пользователя: " + err.Error(),
			})
		}

		fullName, ok := userInfo["full_name"].(string)
		if !ok || fullName == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Полное имя пользователя отсутствует или имеет неверный формат",
			})
		}

		phoneNumber, ok := userInfo["phone_number"].(string)
		fmt.Println(phoneNumber)
		if !ok || phoneNumber == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при получении номера телефона",
			})
		}

		err = db.Table("bookings").
			Joins("INNER JOIN clients ON clients.client_id = bookings.client_id").
			Where("bookings.booker_full_name = ? AND clients.phone_number = ?", fullName, phoneNumber).
			Where("bookings.booking_time >= ?", time.Now().Add(-1*time.Hour)).
			Find(&bookingsList).Error

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка при получении бронирований: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"data": bookingsList,
		})
	}
}
