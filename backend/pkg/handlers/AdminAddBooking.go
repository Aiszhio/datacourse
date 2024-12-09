package handlers

import (
	"errors"
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type AdminBooking struct {
	BookingID      int       `gorm:"column:booking_id;primaryKey;autoIncrement" json:"id"`
	BookingType    string    `gorm:"column:booking_type;not null" json:"booking_type"`
	OrderID        int       `gorm:"column:order_id;not null;autoIncrement" json:"order_id"`
	BookingTime    time.Time `gorm:"column:booking_time;not null" json:"booking_time"`
	BookerFullName string    `gorm:"column:booker_full_name;not null" json:"booker_full_name"`
	PhoneNumber    string    `gorm:"column:phone_number;not null" json:"phone_number"`
}

func AdminCreateBooking(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var adminBooking AdminBooking

		if err := c.BodyParser(&adminBooking); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Неверный формат данных: " + err.Error(),
			})
		}

		var clientID int

		if err := db.Table("clients").Where("phone_number = ?", adminBooking.PhoneNumber).Pluck("client_id", &clientID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Не удалось найти пользователя с таким номером телефона",
			})
		}

		if clientID == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Клиент с таким номером телефона не найден",
			})
		}

		booking := db2.Booking{
			BookingID:      adminBooking.BookingID,
			BookingType:    adminBooking.BookingType,
			OrderID:        adminBooking.OrderID,
			BookingTime:    adminBooking.BookingTime,
			BookerFullName: adminBooking.BookerFullName,
			ClientID:       clientID,
		}

		err := AdminValidateBooking(db, booking)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var client db2.Client
		if err = db.Table("clients").
			Where("full_name = ?", booking.BookerFullName).
			First(&client).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Клиент с именем '%s' не найден", booking.BookerFullName),
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при поиске клиента: " + err.Error(),
			})
		}

		tx := db.Begin()
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Не удалось начать транзакцию: " + tx.Error.Error(),
			})
		}

		if err = tx.Create(&booking).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Не удалось создать бронирование: " + err.Error(),
			})
		}

		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Не удалось зафиксировать транзакцию: " + err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Бронь успешно создана!",
		})
	}
}

func AdminValidateBooking(db *gorm.DB, booking db2.Booking) error {
	var userID int

	if err := db.Table("clients").Select("client_id").
		Where("full_name = ?", booking.BookerFullName).
		Scan(&userID).Error; err != nil {
		return fmt.Errorf("Не удалось получить идентификатор пользователя: %v", err)
	}

	if userID == 0 {
		return fmt.Errorf("Пользователь с идентификатором '%s' не найден", booking.BookerFullName)
	}

	err := ValidateBookingCommon(db, booking, userID)
	if err != nil {
		return err
	}

	return nil
}
