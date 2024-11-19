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

type ClientOrder struct {
	ServiceID int       `json:"service"`
	OrderDate time.Time `json:"orderDate"`
}

func CreateOrder(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var order ClientOrder
		var orderToSave database.Order

		if err := c.BodyParser(&order); err != nil {
			fmt.Println(order)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request payload",
			})
		}

		if order.ServiceID == 0 || order.OrderDate.IsZero() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Service ID and Order Date are required",
			})
		}

		clientMap, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		clientID, ok := clientMap["client_id"]
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid client ID",
			})
		}

		var clientIdInt int

		switch v := clientID.(type) {
		case int:
			clientIdInt = v
		case float64:
			clientIdInt = int(v)
		case string:
			clientIdInt, err = strconv.Atoi(v)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error parsing client ID",
				})
			}
		default:
			fmt.Println("Error: client_id is not a recognized type (int, float64, or string)")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving client ID",
			})
		}

		orderToSave.ClientID = clientIdInt

		err = db.Table("services").
			Where("service_id = ?", order.ServiceID).Pluck("name", &orderToSave.ServiceName).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving service name",
			})
		}

		orderToSave.OrderDate = order.OrderDate

		orderToSave.ReceiptDate, err = GetReceiptDate(orderToSave.ServiceName, order.OrderDate)
		fmt.Printf("Service Name: %s, Order Date: %v\n", orderToSave.ServiceName, order.OrderDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fmt.Printf("Calculated Receipt Date: %v\n", orderToSave.ReceiptDate)

		orderInfo, err := ValidateOrder(db, order, orderToSave.ReceiptDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		workerIDInterface, ok := orderInfo["worker_id"]
		if !ok {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "No available workers for the selected date",
			})
		}
		workerID, ok := workerIDInterface.(int)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid worker ID",
			})
		}
		orderToSave.EmployeeID = workerID

		if err = db.Create(&orderToSave).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error creating order",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"order": orderToSave,
		})
	}
}
