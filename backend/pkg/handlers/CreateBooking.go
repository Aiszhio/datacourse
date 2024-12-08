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

type UserBooking struct {
	BookingType string    `gorm:"column:booking_type;not null" json:"booking_type"`
	BookingTime time.Time `gorm:"column:booking_time;not null" json:"booking_time"`
	BookerName  string    `gorm:"column:booker_full_name;not null" json:"booker_full_name"`
}

type BookingOrder struct {
	BookingID int `gorm:"booking_id"`
}

func CreateBooking(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var booking UserBooking
		var bookingToOrder database.BookingToOrder

		if err := c.BodyParser(&booking); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		workerID, err := ValidateOrder(db, booking, client)
		if err != nil {
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

		err = db.Table("bookings").Create(&booking).Error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var workerName string
		err = db.Table("employees").Where("employee_id = ?", workerID).
			Select("full_name").Scan(&workerName).Error
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		bookingToOrder.EmployeeID = workerID
		bookingToOrder.OrderDate = bookingTime

		err = PrepareOrder(db, client, bookingToOrder)
		if err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "Successfully created booking",
			"worker": workerName,
		})
	}
}

func PrepareOrder(db *gorm.DB, client *redis.Client, orders database.BookingToOrder) error {
	clientInfo, err := Redis.GetMultipleKey(client, "UserInfo")
	if err != nil {
		fmt.Println("Could not get user info")
		return err
	}

	var clientID int

	switch v := clientInfo["client_id"].(type) {
	case string:
		clientID, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	case int:
		clientID = v
	case float64:
		clientID = int(v)
	default:
		fmt.Println("Something is wrong with client info")
		return err
	}

	orders.ClientID = clientID

	var bookingOrder BookingOrder

	err = db.Table("bookings").Select("booking_id").Where("booking_time = ?", orders.OrderDate).
		First(&bookingOrder).Error
	if err != nil {
		fmt.Println(err)
		return err
	}

	orders.BookingID = bookingOrder.BookingID

	orders.OrderEnd = orders.OrderDate.Add(time.Hour)

	err = db.Table("booking_to_orders").Create(&orders).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
