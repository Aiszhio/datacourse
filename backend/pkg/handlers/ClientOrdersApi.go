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

type OrderWithEmployeeName struct {
	OrderID      int       `json:"order_id"`
	ClientID     int       `json:"client_id"`
	EmployeeName string    `json:"employee_name"`
	ServiceName  string    `json:"service_name"`
	OrderDate    time.Time `json:"order_date"`
	ReceiptDate  time.Time `json:"receipt_date"`
}

func ClientOrdersApi(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
		fmt.Println("Retrieved data from Redis:", userInfo)
		if err != nil {
			fmt.Println("Error retrieving user info from Redis:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving user info from Redis",
			})
		}

		var clientID int
		switch v := userInfo["client_id"].(type) {
		case int:
			clientID = v
		case float64:
			clientID = int(v)
		case string:
			clientID, err = strconv.Atoi(v)
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

		fmt.Println("Client ID retrieved from Redis:", clientID)

		var orders []OrderWithEmployeeName

		if err = db.Table("orders").
			Select("orders.order_id, orders.client_id, employees.full_name as employee_name, orders.service_name, orders.order_date, orders.receipt_date").
			Joins("left join employees on orders.employee_id = employees.employee_id").
			Where("orders.client_id = ?", clientID).
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
