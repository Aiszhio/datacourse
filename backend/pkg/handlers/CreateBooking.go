package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func CreateBooking(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var booking database.Booking

		if err := c.BodyParser(&booking); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		bookingTime := booking.BookingTime
		moscowLocation, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Println("Error loading Moscow timezone:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load Moscow timezone",
			})
		}

		bookingTime = bookingTime.In(moscowLocation)

		err = ValidateBooking(db, booking, client)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при работе с номером клиента",
			})
		}

		clientIDInt, ok := userInfo["client_id"]
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка в конвертации ID",
			})
		}

		var clientID int

		switch v := clientIDInt.(type) {
		case string:
			clientID, err = strconv.Atoi(v)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Ошибка в конвертации ID из строчки",
				})
			}
		case int:
			clientID = v
		case float64:
			clientID = int(v)
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка в конвертации данных",
			})
		}

		booking.ClientID = clientID

		err = db.Table("bookings").Create(&booking).Error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "Successfully created booking",
		})
	}
}
