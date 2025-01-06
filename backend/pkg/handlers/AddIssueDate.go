package handlers

import (
	"errors"
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/Aiszhio/datacourse.git/pkg/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type IssueRequest struct {
	IssueDate string `json:"issueDate"`
}

func IssueOrder(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		fmt.Println("Got ID param:", idStr)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Error converting ID:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Некорректный ID: " + err.Error(),
			})
		}

		var req IssueRequest
		if err = c.BodyParser(&req); err != nil {
			fmt.Println("Error parsing body:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга JSON: " + err.Error(),
			})
		}
		fmt.Printf("Parsed IssueRequest: %#v\n", req)

		if req.IssueDate == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Поле issueDate не может быть пустым",
			})
		}

		t, err := time.Parse(time.RFC3339, req.IssueDate)
		if err != nil {
			fmt.Println("Time parse error:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Ошибка парсинга issueDate: %v", err),
			})
		}
		fmt.Println("Parsed time:", t)

		if err = utils.CheckWorkingHours(t); err != nil {
			fmt.Println("CheckWorkingHours failed:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		var order database.Order
		if err = db.First(&order, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Заказ не найден",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при поиске заказа: " + err.Error(),
			})
		}
		fmt.Printf("Found order: ID=%d, receiveDate=%v\n", order.OrderID, order.ReceiptDate)

		if t.Before(order.ReceiptDate) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Дата выдачи не может быть раньше, чем дата завершения (receive_date)",
			})
		}

		order.IssueDate = t
		if err = db.Table("orders").Save(&order).Error; err != nil {
			fmt.Println("Save error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при сохранении заказа: " + err.Error(),
			})
		}

		fmt.Println("Заказ успешно обновлён:", order.OrderID)
		return c.JSON(fiber.Map{
			"message": "Заказ успешно выдан",
		})
	}
}
