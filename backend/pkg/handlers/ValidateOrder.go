package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func ValidateOrder(db *gorm.DB, order UserBooking, client *redis.Client) (int, error) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return 0, fmt.Errorf("Не удалось загрузить часовой пояс: %v", err)
	}

	bookingTime := order.BookingTime.In(loc)
	now := time.Now().In(loc)

	err = CheckOtherOrders(db, order.BookingTime, client)
	if err != nil {
		return 0, err
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

func CheckOtherOrders(db *gorm.DB, orderTime time.Time, client *redis.Client) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return err
	}

	userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
	if err != nil {
		return err
	}

	var userID int

	switch v := userInfo["client_id"].(type) {
	case string:
		userID, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		fmt.Println("Error in recognizing type of user id")
		return fmt.Errorf("Invalid client ID")
	}

	var count int64

	err = db.Table("booking_to_orders").
		Where("client_id = ?", userID).
		Where("order_date < ? AND order_end > ?", orderTime.In(loc).Add(time.Hour), orderTime.In(loc)).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("Ошибка при проверке других заказов: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("У вас уже есть бронирование на это время")
	}

	return nil
}
