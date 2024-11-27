package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

func AdminCreateBooking(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var booking UserBooking

		if err := c.BodyParser(&booking); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		workerID, err := AdminValidateBooking(db, booking)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var clientID int
		if err = db.Table("clients").
			Select("client_id").
			Where("full_name = ?", booking.BookerName).
			Scan(&clientID).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tx := db.Begin()

		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": tx.Error.Error(),
			})
		}

		bookingToSave := db2.Booking{
			BookingType:    booking.BookingType,
			BookingTime:    booking.BookingTime,
			BookerFullName: booking.BookerName,
		}

		if err = tx.Table("bookings").Omit("order_id").Create(&bookingToSave).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		bookingToOrder := db2.BookingToOrder{
			BookingID:  bookingToSave.BookingID,
			EmployeeID: workerID,
			ClientID:   clientID,
			OrderDate:  booking.BookingTime,
			OrderEnd:   booking.BookingTime.Add(time.Hour),
		}

		if err = tx.Table("booking_to_orders").Create(&bookingToOrder).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tx.Commit()

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Бронь успешно создана!",
		})
	}
}

func AdminValidateBooking(db *gorm.DB, booking UserBooking) (int, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return 0, fmt.Errorf("Не удалось загрузить часовой пояс: %v", err)
	}

	bookingTime := booking.BookingTime.In(loc)
	now := time.Now().In(loc)

	var userID int
	if err = db.Table("clients").Select("client_id").
		Where("full_name = ?", booking.BookerName).
		Scan(&userID).Error; err != nil {
		return 0, err
	}

	if userID == 0 {
		return 0, fmt.Errorf("Пользователь с именем '%s' не найден", booking.BookerName)
	}

	err = CheckOtherBookings(db, booking, userID)
	if err != nil {
		return 0, fmt.Errorf("У пользователя уже есть бронирование на это время")
	}

	if bookingTime.Before(now) {
		return 0, fmt.Errorf("Время бронирования не может быть в прошлом")
	}

	openingHour := 9
	closingHour := 19

	year, month, day := bookingTime.Date()

	openingTime := time.Date(year, month, day, openingHour, 0, 0, 0, loc)
	closingTime := time.Date(year, month, day, closingHour, 0, 0, 0, loc)

	if bookingTime.Before(openingTime) || !bookingTime.Before(closingTime) {
		return 0, fmt.Errorf("Бронирование возможно только с 9:00 до 19:00")
	}

	orderEnd := bookingTime.Add(time.Hour)

	if !orderEnd.Before(closingTime) {
		return 0, fmt.Errorf("Время окончания бронирования не может быть позже 19:00")
	}

	workerID, err := SetWorker(db, bookingTime)
	if err != nil {
		return 0, fmt.Errorf("Ошибка при назначении работника: %v", err)
	}

	if workerID == 0 {
		return 0, fmt.Errorf("Нет доступных работников")
	}

	return workerID, nil
}

func CheckOtherBookings(db *gorm.DB, booking UserBooking, userID int) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	var count int64

	err = db.Table("booking_to_orders").
		Where("client_id = ?", userID).
		Where("order_date < ? AND order_end > ?", booking.BookingTime.In(loc).Add(time.Hour), booking.BookingTime.In(loc)).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("Ошибка при проверке других заказов: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("У вас уже есть бронирование на это время")
	}

	return nil
}
