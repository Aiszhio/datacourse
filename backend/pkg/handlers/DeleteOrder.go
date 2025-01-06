package handlers

import (
	"errors"
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteOrder(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("orderId")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка конвертации id",
			})
		}

		var order db2.Order
		err = db.Table("orders").Where("order_id = ?", id).First(&order).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Такого заказа не существует",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Ошибка при поиске заказа: %v", err),
			})
		}

		if !order.IssueDate.IsZero() && order.IssueDate.Before(time.Now()) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Нельзя удалить прошедший заказ",
			})
		}

		if err = db.Table("orders").Where("order_id = ?", id).Delete(&order).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Ошибка при удалении заказа: %v", err),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Заказ успешно удалён",
		})
	}
}
