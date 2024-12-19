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

// Предполагается, что у вас есть структура Order, которая отражает таблицу orders
type Order struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ClientID    uint      `json:"client_id"`
	ServiceName string    `json:"service_name"`
	OrderDate   time.Time `json:"order_date"`
	// Другие поля заказа
}

// Функция DeleteClient обрабатывает удаление клиента
func DeleteClient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем параметр "id" из URL
		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Требуется идентификатор клиента",
			})
		}

		// Конвертируем "id" из строки в целое число
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка конвертации id клиента",
			})
		}

		// Получаем существующего клиента из базы данных
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

		// Получаем текущее время в UTC
		nowUtc := time.Now().UTC()

		// Проверяем наличие не закрытых заказов у клиента
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

		// Удаляем клиента из базы данных
		if err = db.Delete(&existingClient).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при удалении клиента",
			})
		}

		// Возвращаем успешный ответ
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Клиент успешно удален!",
		})
	}
}
