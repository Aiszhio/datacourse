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
		fmt.Printf("Received workerIDStr: %s\n", workerIDStr)

		if workerIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Employee ID is required.",
			})
		}

		workerID, err := strconv.Atoi(workerIDStr)
		if err != nil {
			fmt.Printf("Error converting workerIDStr to int: %v\n", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid Employee ID.",
			})
		}
		fmt.Printf("Parsed workerID: %d\n", workerID)

		var employee db2.Employee
		if err = db.First(&employee, workerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Сотрудник не найден",
				})
			}
			fmt.Printf("Error retrieving employee: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при поиске сотрудника",
			})
		}

		if employee.Status == "Уволен" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Сотрудник уже уволен",
			})
		}

		nowUtc := time.Now().UTC()
		fmt.Printf("Current time in UTC: %s\n", nowUtc.Format(time.RFC3339))

		var upcomingOrders int64
		if err = db.Table("orders").
			Where("employee_id = ? AND order_date > ?", workerID, nowUtc).
			Count(&upcomingOrders).Error; err != nil {
			fmt.Printf("Error checking upcoming orders: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка при проверке предстоящих заказов.",
			})
		}

		fmt.Printf("Upcoming orders for employee %d: %d\n", workerID, upcomingOrders)

		if upcomingOrders > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Сотрудника нельзя уволить, так как у него есть предстоящие заказы.",
			})
		}

		if err = db.Model(&employee).Update("status", "Уволен").Error; err != nil {
			fmt.Printf("Error updating employee status: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка при обновлении статуса сотрудника.",
			})
		}

		fmt.Printf("Employee %d fired successfully.\n", workerID)

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Printf("Error loading Moscow location: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка при загрузке местоположения.",
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
