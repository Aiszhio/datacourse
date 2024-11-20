package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type OrderWithClientNames struct {
	OrderID     int       `json:"order_id"`
	ClientName  string    `json:"client_name"`
	ServiceName string    `json:"service_name"`
	OrderDate   time.Time `json:"order_date"`
	ReceiptDate time.Time `json:"receipt_date"`
}

func WorkerOrdersApi(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var orders []OrderWithClientNames

		workerMap, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			fmt.Println("Error getting user info")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "cannot get user info",
			})
		}

		workerID := workerMap["employee_id"]

		switch v := workerID.(type) {
		case int:
			workerID = v
		case string:
			workerID, err = strconv.Atoi(v)
			if err != nil {
				fmt.Println("Error converting workerID to int")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"Error": "Error converting workerID to int",
				})
			}
		case float64:
			workerID = int(v)
		default:
			fmt.Println("Error: employee_id is not a recognized type (int, float64, or string)")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving employee ID",
			})
		}

		if err = db.Table("orders").
			Select("orders.order_id, clients.full_name as client_name, "+
				"orders.service_name, orders.order_date, orders.receipt_date").
			Joins("left join employees on orders.employee_id = employees.employee_id").
			Joins("left join clients on orders.client_id = clients.client_id").
			Where("orders.employee_id = ?", workerID).
			Find(&orders).Error; err != nil {
			fmt.Println("Error retrieving orders from database:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving orders from database",
			})
		}

		fmt.Println("Orders successfully fetched: ", orders)
		return c.JSON(fiber.Map{
			"data": orders,
		})
	}
}
