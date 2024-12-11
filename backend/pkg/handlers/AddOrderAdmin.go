package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type ClientOrder struct {
	ClientName  string    `json:"clientName"`
	ServiceName string    `json:"service"`
	OrderDate   time.Time `json:"orderDate"`
}

func AddOrder(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var clientOrder ClientOrder
		var order db2.Order

		if err := c.BodyParser(&clientOrder); err != nil {
			fmt.Println("Ошибка парсинга данных:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга данных",
			})
		}

		order.ServiceName = clientOrder.ServiceName
		order.OrderDate = clientOrder.OrderDate
		order.ReceiptDate = clientOrder.OrderDate.Add(time.Hour)

		var client db2.Client
		if err := db.Table("clients").Where("full_name = ?", clientOrder.ClientName).First(&client).Error; err != nil {
			fmt.Println("Ошибка поиска клиента:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Пользователь с таким ФИО не найден",
			})
		}
		order.ClientID = client.ClientID

		employeeID, err := SetWorker(db, order.OrderDate)
		if err != nil {
			fmt.Println("Ошибка функции SetWorker:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Не удалось найти свободного работника",
			})
		}

		if employeeID == 0 {
			fmt.Println("Все работники заняты")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "На данный момент все работники заняты",
			})
		}

		order.EmployeeID = employeeID

		var service db2.Service
		if err = db.Table("services").Where("name = ?", clientOrder.ServiceName).First(&service).Error; err != nil {
			fmt.Println("Ошибка поиска услуги:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Такая услуга не найдена",
			})
		}

		orderContent := db2.OrderContent{
			ServiceID: service.ServiceID,
		}

		if err = ValidateOrderEquipment(order, db); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tx := db.Begin()
		if tx.Error != nil {
			fmt.Println("Ошибка начала транзакции:", tx.Error)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка при начале транзакции",
			})
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				fmt.Println("Паника при транзакции:", r)
			}
		}()

		if err = tx.Create(&order).Error; err != nil {
			tx.Rollback()
			fmt.Println("Ошибка создания заказа:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка при записи заказа в БД",
			})
		}

		orderContent.OrderID = order.OrderID
		if err = tx.Create(&orderContent).Error; err != nil {
			tx.Rollback()
			fmt.Println("Ошибка создания содержимого заказа:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка при записи содержимого заказа в БД",
			})
		}

		if err = tx.Commit().Error; err != nil {
			fmt.Println("Ошибка коммита транзакции:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка при завершении транзакции",
			})
		}

		fmt.Println(order.OrderDate, order.ReceiptDate)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Заказ успешно сохранен",
			"order":   order,
		})
	}
}
