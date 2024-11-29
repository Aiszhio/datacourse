package handlers

import (
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"sort"
	"time"
)

type OrdersWithNames struct {
	OrderID      int       `json:"id"`
	ClientName   string    `json:"clientName"`
	EmployeeName string    `json:"employeeName"`
	ServiceName  string    `json:"service"`
	OrderDate    time.Time `json:"orderDate"`
	ReceiptDate  time.Time `json:"receiveDate"`
}

func GetAdminOrders(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var order []database.Order

		if err := db.Table("orders").Select("order_id, client_id, employee_id, service_name," +
			"order_date, receipt_date").Find(&order).Error; err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err,
			})
		}

		sendOrders := FillOrdersData(order, db)

		if sendOrders == nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error occurred while filling data for admin orders",
			})
		}

		sort.Slice(sendOrders, func(i, j int) bool {
			return sendOrders[i].OrderID < sendOrders[j].OrderID
		})

		return c.JSON(fiber.Map{
			"orders": sendOrders,
		})
	}
}

func FillOrdersData(order []database.Order, db *gorm.DB) []OrdersWithNames {
	var sendOrders []OrdersWithNames

	for i := 0; i < len(order); i++ {

		tempClientName, err := GetClientName(order[i].ClientID, db)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		tempEmployeeName, err := GetWorkerName(order[i].EmployeeID, db)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		sendOrders = append(sendOrders, OrdersWithNames{
			OrderID:      order[i].OrderID,
			ClientName:   tempClientName,
			EmployeeName: tempEmployeeName,
			ServiceName:  order[i].ServiceName,
			OrderDate:    order[i].OrderDate,
			ReceiptDate:  order[i].ReceiptDate,
		})
	}

	if len(sendOrders) == 0 {
		fmt.Println("Error in searching all orders")
		return nil
	}

	return sendOrders
}
