package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ValidateBookingCommon(db *gorm.DB, booking db.Booking, userID int) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return fmt.Errorf("Не удалось загрузить часовой пояс: %v", err)
	}

	bookingTime := booking.BookingTime.In(loc)
	now := time.Now().In(loc)

	if bookingTime.Before(now) {
		fmt.Printf("Booking Time: %v, Now: %v\n", bookingTime, now)
		return fmt.Errorf("Время бронирования не может быть в прошлом")
	}

	openingHour := 9
	closingHour := 19

	year, month, day := bookingTime.Date()

	openingTime := time.Date(year, month, day, openingHour, 0, 0, 0, loc)
	closingTime := time.Date(year, month, day, closingHour, 0, 0, 0, loc)

	if bookingTime.Before(openingTime) || !bookingTime.Before(closingTime) {
		return fmt.Errorf("Бронирование возможно только с 9:00 до 19:00")
	}

	orderEnd := bookingTime.Add(time.Hour)
	if !orderEnd.Before(closingTime) {
		return fmt.Errorf("Время окончания бронирования не может быть позже 19:00")
	}

	err = CheckOtherBookingsCommon(db, booking, userID)
	if err != nil {
		return err
	}

	return nil
}

func CheckOtherBookingsCommon(db *gorm.DB, booking db.Booking, userID int) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	var count int64

	bookingStart := booking.BookingTime.In(loc)
	bookingEnd := bookingStart.Add(time.Hour)

	err = db.Table("bookings").
		Where("booker_full_name = ?", booking.BookerFullName).
		Where("booking_time < ? AND booking_time + interval '1 hour' > ?", bookingEnd, bookingStart).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("Ошибка при проверке других бронирований: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("У вас уже есть бронирование на это время")
	}

	return nil
}

func ValidateBooking(db *gorm.DB, booking db.Booking, client *redis.Client) error {
	userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
	if err != nil {
		return fmt.Errorf("Не удалось получить информацию о пользователе: %v", err)
	}

	var userID int

	switch v := userInfo["client_id"].(type) {
	case string:
		userID, err = strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("Неверный формат client_id: %v", err)
		}
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		fmt.Println("Ошибка при распознавании типа user_id")
		return fmt.Errorf("Неверный client ID")
	}

	return ValidateBookingCommon(db, booking, userID)
}
