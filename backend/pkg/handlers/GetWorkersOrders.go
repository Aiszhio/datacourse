package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

// WorkerOrders представляет структуру данных заказа с добавленным полем IssueDate
type WorkerOrders struct {
	OrderID     int        `json:"id" gorm:"column:order_id"`
	ClientName  string     `json:"clientName" gorm:"column:client_name"`
	ServiceName string     `json:"serviceName" gorm:"column:service_name"`
	OrderDate   time.Time  `json:"orderDate" gorm:"column:order_date"`
	ReceiptDate time.Time  `json:"receiptDate" gorm:"column:receipt_date"`
	IssueDate   *time.Time `json:"issueDate" gorm:"column:issue_date"` // Новое поле
}

// WorkerOrderResponse заменяет RemainingTime на IssueDate
type WorkerOrderResponse struct {
	OrderID         int        `json:"id"`
	ClientName      string     `json:"clientName"`
	ServiceName     string     `json:"serviceName"`
	OrderDate       time.Time  `json:"orderDate"`
	AdjustedReceipt time.Time  `json:"adjustedReceiptDate"`
	IssueDate       *time.Time `json:"issueDate,omitempty"`
}

type OrdersMapResponse struct {
	PastOrders   []WorkerOrderResponse `json:"pastOrders"`
	FutureOrders []WorkerOrderResponse `json:"futureOrders"`
}

func GetWorkerOrders(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var orders []WorkerOrders

		userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			log.Printf("Ошибка при получении данных о пользователе: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при получении данных о пользователе",
			})
		}

		var employeeId int

		switch v := userInfo["employee_id"].(type) {
		case string:
			employeeId, err = strconv.Atoi(v)
			if err != nil {
				log.Printf("Ошибка форматирования employeeId: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Ошибка форматирования данных",
				})
			}
		case int:
			employeeId = v
		case float64:
			employeeId = int(v)
		default:
			log.Printf("Ошибка приведения типов для employeeId: %v", v)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка приведения типов",
			})
		}

		log.Printf("Получение заказов для employeeId: %d", employeeId)

		if err = db.Table("orders").
			Select("orders.order_id, clients.full_name as client_name, services.name as service_name, orders.order_date, orders.receipt_date, orders.issue_date").
			Joins("JOIN clients ON orders.client_id = clients.client_id").
			Joins("JOIN services ON orders.service_name = services.name").
			Where("orders.employee_id = ? AND orders.service_name IS NOT NULL AND orders.service_name <> ''", employeeId).
			Scan(&orders).Error; err != nil {
			log.Printf("Ошибка при выполнении запроса к базе данных: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при получении заказов",
			})
		}

		log.Printf("Получено заказов: %d", len(orders))

		for _, order := range orders {
			issueDateStr := "NULL"
			if order.IssueDate != nil {
				issueDateStr = order.IssueDate.Format(time.RFC3339)
			}
			log.Printf("Заказ ID: %d, ClientName: '%s', ServiceName: '%s', OrderDate: %s, ReceiptDate: %s, IssueDate: %s",
				order.OrderID, order.ClientName, order.ServiceName, order.OrderDate.Format(time.RFC3339), order.ReceiptDate.Format(time.RFC3339), issueDateStr)
		}

		ordersMap, err := GetOrders(db, orders)
		if err != nil {
			log.Printf("Ошибка при разделении заказов: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		log.Printf("Прошедших заказов: %d, Будущих заказов: %d", len(ordersMap.PastOrders), len(ordersMap.FutureOrders))

		return c.JSON(fiber.Map{
			"message": "Заказы успешно загружены",
			"orders":  ordersMap,
		})
	}
}

func GetOrders(db *gorm.DB, orders []WorkerOrders) (OrdersMapResponse, error) {
	pastOrders := make([]WorkerOrderResponse, 0)
	futureOrders := make([]WorkerOrderResponse, 0)
	now := time.Now()

	serviceNames := make([]string, 0, len(orders))
	for _, order := range orders {
		serviceNames = append(serviceNames, order.ServiceName)
	}

	var services []database.Service
	if err := db.Table("services").Where("name IN ?", serviceNames).Find(&services).Error; err != nil {
		return OrdersMapResponse{}, fmt.Errorf("ошибка при получении данных об услугах: %w", err)
	}

	serviceMap := make(map[string]database.Service)
	for _, service := range services {
		serviceMap[service.Name] = service
	}

	for _, order := range orders {
		service, exists := serviceMap[order.ServiceName]
		if !exists {
			log.Printf("Услуга '%s' для заказа ID %d не найдена", order.ServiceName, order.OrderID)
			continue
		}

		adjustedReceiptDate := order.ReceiptDate
		log.Printf("Заказ ID %d: adjustedReceiptDate = %s", order.OrderID, adjustedReceiptDate.Format(time.RFC3339))

		if order.IssueDate != nil {
			classificationTime := *order.IssueDate
			isPast := classificationTime.Before(now)
			log.Printf("Заказ ID %d: IssueDate = %s, isPast = %v", order.OrderID, classificationTime.Format(time.RFC3339), isPast)

			if isPast {
				pastOrders = append(pastOrders, WorkerOrderResponse{
					OrderID:         order.OrderID,
					ClientName:      order.ClientName,
					ServiceName:     order.ServiceName,
					OrderDate:       order.OrderDate,
					AdjustedReceipt: adjustedReceiptDate,
					IssueDate:       order.IssueDate,
				})
			} else {
				futureOrders = append(futureOrders, WorkerOrderResponse{
					OrderID:         order.OrderID,
					ClientName:      order.ClientName,
					ServiceName:     order.ServiceName,
					OrderDate:       order.OrderDate,
					AdjustedReceipt: adjustedReceiptDate,
					IssueDate:       order.IssueDate,
				})
				log.Printf("Заказ ID %d классифицирован как будущий, IssueDate = %v", order.OrderID, order.IssueDate)
			}
		} else {
			remainingTime := adjustedReceiptDate.AddDate(0, 0, service.Days).Format(time.RFC3339)

			parsedRemainingTime, err := time.Parse(time.RFC3339, remainingTime)
			if err != nil {
				log.Printf("Ошибка парсинга remainingTime для заказа ID %d: %v", order.OrderID, err)
				continue
			}

			isPast := parsedRemainingTime.Before(now)
			log.Printf("Заказ ID %d: RemainingTime = %s, isPast = %v", order.OrderID, remainingTime, isPast)

			if isPast {
				pastOrders = append(pastOrders, WorkerOrderResponse{
					OrderID:         order.OrderID,
					ClientName:      order.ClientName,
					ServiceName:     order.ServiceName,
					OrderDate:       order.OrderDate,
					AdjustedReceipt: adjustedReceiptDate,
					IssueDate:       order.IssueDate, // nil
				})
			} else {
				futureOrders = append(futureOrders, WorkerOrderResponse{
					OrderID:         order.OrderID,
					ClientName:      order.ClientName,
					ServiceName:     order.ServiceName,
					OrderDate:       order.OrderDate,
					AdjustedReceipt: adjustedReceiptDate,
					IssueDate:       order.IssueDate, // nil
				})
				log.Printf("Заказ ID %d классифицирован как будущий, RemainingTime = %s", order.OrderID, remainingTime)
			}
		}
	}

	log.Printf("Прошедших заказов: %d, Будущих заказов: %d", len(pastOrders), len(futureOrders))

	return OrdersMapResponse{
		PastOrders:   pastOrders,
		FutureOrders: futureOrders,
	}, nil
}
