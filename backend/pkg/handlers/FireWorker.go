package handlers

import (
	"errors"
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func FireWorker(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		workerIDStr := c.Params("id")
		if workerIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Отсутствует идентификатор сотрудника.",
			})
		}

		workerID, err := strconv.Atoi(workerIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неверный формат идентификатора сотрудника.",
			})
		}

		var employee db2.Employee
		if err = db.First(&employee, workerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Сотрудник не найден.",
				})
			}
			fmt.Printf("Ошибка при получении сотрудника из БД: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при поиске сотрудника в базе.",
			})
		}

		if employee.Status == "Уволен" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Сотрудник уже уволен.",
			})
		}

		nowUtc := time.Now().UTC()

		var upcomingOrders int64
		// Допустим, что services.days — это целочисленный тип, совместимый с int.
		// Если это bigint, то делаем явное приведение к int: services.days::int
		if err = db.Table("orders").
			Joins("INNER JOIN services ON orders.service_name = services.name").
			Where("orders.employee_id = ? AND (orders.order_date + make_interval(0,0,0,services.days::int,0,0,0)) > ?", workerID, nowUtc).
			Count(&upcomingOrders).Error; err != nil {
			fmt.Printf("Ошибка при проверке предстоящих заказов: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка при проверке предстоящих заказов.",
			})
		}

		if upcomingOrders > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Невозможно уволить сотрудника, так как у него есть предстоящие заказы.",
			})
		}

		if err = db.Model(&employee).Update("status", "Уволен").Error; err != nil {
			fmt.Printf("Ошибка при обновлении статуса сотрудника: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка при обновлении статуса сотрудника.",
			})
		}

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Printf("Ошибка при загрузке временной зоны Москва: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка при настройке временной зоны.",
			})
		}

		hireDateMsk := employee.HireDate.In(loc).Format(time.RFC3339)
		birthDateMsk := employee.BirthDate.In(loc).Format(time.RFC3339)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Сотрудник успешно уволен.",
			"worker": fiber.Map{
				"employee_id":   employee.EmployeeID,
				"full_name":     employee.FullName,
				"position":      employee.Position,
				"hire_date":     hireDateMsk,
				"birth_date":    birthDateMsk,
				"passport_data": employee.PassportData,
				"phone_number":  employee.PhoneNumber,
				"status":        employee.Status,
			},
		})
	}
}
