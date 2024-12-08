package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
	"time"
)

type ClientOrder struct {
	ClientName  string    `json:"clientName"`
	ServiceName string    `json:"service"`
	OrderDate   time.Time `json:"orderDate"`
}

func AddOrder(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var order ClientOrder
		var orderFromBooking db2.BookingToOrder
		var orderToService db2.OrderContent

		if err := c.BodyParser(&order); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга данных",
			})
		}

		if order.ClientName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Имя клиента не может быть пустым",
			})
		}

		if order.ServiceName == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Название услуги не может быть пустым",
			})
		}

		tx := db.Begin()
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при запуске транзакции",
			})
		}
		defer tx.Rollback()

		order.ClientName = strings.TrimSpace(order.ClientName)
		order.ServiceName = strings.TrimSpace(order.ServiceName)

		var client db2.Client
		if err := tx.Table("clients").Where("full_name = ?", order.ClientName).First(&client).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Такого клиента не существует",
			})
		}

		var service db2.Service
		if err := tx.Table("services").Select("service_id").Where("name = ?", order.ServiceName).First(&service).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Такой услуги не существует",
			})
		}

		moscowLocation, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при получении часового пояса для Москвы",
			})
		}

		orderDate := order.OrderDate.In(moscowLocation)

		orderDateStr := orderDate.Format(time.RFC3339)

		fmt.Println("Parsed Order Date:", orderDateStr)

		if orderDate.Hour() < 9 || orderDate.Hour() >= 19 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Время заказа должно быть между 9:00 и 19:00 МСК",
			})
		}

		if err = tx.Table("booking_to_orders as bto").
			Select("bto.client_id, bto.order_date, c.full_name as client_name").
			Joins("inner join clients as c on bto.client_id = c.client_id").
			Where("bto.client_id = ? and bto.order_date < ?", client.ClientID, order.OrderDate).
			First(&orderFromBooking).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Такой брони не существует",
			})
		}

		receiptDate := orderFromBooking.OrderDate.In(moscowLocation).Add(time.Hour)

		orderToSave := db2.Order{
			ClientID:    orderFromBooking.ClientID,
			EmployeeID:  orderFromBooking.EmployeeID,
			ServiceName: order.ServiceName,
			OrderDate:   orderDate,
			ReceiptDate: receiptDate,
		}

		if err = tx.Table("orders").Create(&orderToSave).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка сохранения заказа",
			})
		}

		orderToService.ServiceID = service.ServiceID

		if err = tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при подтверждении транзакции",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Заказ успешно сохранен",
		})
	}
}
