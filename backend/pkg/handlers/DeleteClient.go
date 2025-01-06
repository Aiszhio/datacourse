package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClientID    uint      `json:"client_id"`
	ServiceName string    `json:"service_name"`
	OrderDate   time.Time `json:"order_date"`
}

func DeleteClient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Требуется идентификатор клиента",
			})
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка конвертации id клиента",
			})
		}

		var existingClient db2.Client
		if err = db.First(&existingClient, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Клиент не найден",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при поиске клиента",
			})
		}

		nowUtc := time.Now().UTC()

		var upcomingOrders int64
		err = db.Table("orders").
			Joins("INNER JOIN services ON orders.service_name = services.name").
			Where("orders.client_id = ? AND (orders.order_date + make_interval(0,0,0,services.days::int,0,0,0)) > ?", id, nowUtc).
			Count(&upcomingOrders).Error
		if err != nil {
			fmt.Printf("Ошибка при проверке предстоящих заказов: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при проверке предстоящих заказов.",
			})
		}

		if upcomingOrders > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Данный клиент не может быть удален, у него есть не закрытые заказы",
			})
		}

		var upcomingBookings int64
		if err = db.Table("bookings").Where("client_id = ? AND booking_time > ?", id, nowUtc).Count(&upcomingBookings).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при поиске бронирований",
			})
		}

		if upcomingBookings > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "У пользователя есть активные брони",
			})
		}

		if err = db.Delete(&existingClient).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при удалении клиента",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Клиент успешно удален!",
		})
	}
}
