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

type OrderWithNames struct {
	OrderID      int       `json:"order_id"`
	ClientName   string    `json:"client_name"`
	EmployeeName string    `json:"employee_name"`
	ServiceName  string    `json:"service_name"`
	OrderDate    time.Time `json:"order_date"`
	ReceiptDate  time.Time `json:"receipt_date"`
}

func ClientOrdersApi(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		workerMap, err := Redis.GetMultipleKey(client, "UserInfo")
		fmt.Println("Retrieved data from Redis:", workerMap)
		if err != nil {
			fmt.Println("Error retrieving user info from Redis:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		workerID, ok := workerMap["employee_id"]
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "no employee_id found",
			})
		}

		switch v := workerID.(type) {
		case int:
			workerID = v
		case float64:
			workerID = int(v)
		case string:
			workerID, err = strconv.Atoi(v)
			if err != nil {
				fmt.Println("Error converting client_id from string to integer:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Error converting client ID",
				})
			}
		default:
			fmt.Println("Error: client_id is not a recognized type (int, float64, or string)")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving client ID",
			})
		}

		fmt.Println("Worker ID retrieved from Redis:", workerID)

		var orders []OrderWithNames

		if err = db.Table("orders").
			Select("orders.order_id, employees.full_name as employee_name, clients.full_name as client_name, "+
				"orders.service_name, orders.order_date, orders.receipt_date").
			Joins("left join employees on orders.employee_id = employees.employee_id").
			Joins("left join clients on orders.client_id = clients.client_id").
			Where("orders.employee_id = ?", workerID).
			Find(&orders).Error; err != nil {
			fmt.Println("Error fetching orders from database:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error fetching orders from database",
			})
		}

		fmt.Println("Orders fetched successfully:", orders)
		return c.JSON(orders)
	}
}
