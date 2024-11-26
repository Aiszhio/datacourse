package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func WorkerOrdersApi(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve user info from Redis",
			})
		}

		var userID int

		switch v := userInfo["employee_id"].(type) {
		case string:
			userID, err = strconv.Atoi(v)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid employee_id format",
				})
			}
		case int:
			userID = v
		case float64:
			userID = int(v)
		default:
			fmt.Println("Error in type assertion employee id")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid employee_id type",
			})
		}

		var upcomingOrders []db2.BookingToOrder

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to load time location",
			})
		}
		now := time.Now().In(loc)

		err = db.Table("booking_to_orders").Where("employee_id = ?", userID).Find(&upcomingOrders).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve upcoming orders",
			})
		}

		var expiredOrdersWithNames []OrdersWithNames

		err = db.Table("orders").
			Select("orders.order_id as order_id, clients.full_name as client_name,"+
				" employees.full_name as employee_name, orders.service_name as service_name, orders.order_date, orders.receipt_date").
			Joins("JOIN clients ON clients.client_id = orders.client_id").
			Joins("JOIN employees ON employees.employee_id = orders.employee_id").
			Where("orders.employee_id = ?", userID).Where("orders.order_date < ?", now).
			Find(&expiredOrdersWithNames).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve expired orders with names",
			})
		}

		return c.JSON(fiber.Map{
			"upcoming": upcomingOrders,
			"expired":  expiredOrdersWithNames,
		})
	}
}
