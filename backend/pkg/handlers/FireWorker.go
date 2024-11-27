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
					"message": "Employee not found.",
				})
			}
			fmt.Printf("Error retrieving employee: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error retrieving employee.",
			})
		}

		if employee.Status == "Уволен" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Employee is already fired.",
			})
		}

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Printf("Error loading location: %v\n", err) // Для отладки
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error loading location.",
			})
		}
		now := time.Now().In(loc)
		fmt.Printf("Current time in Moscow: %s\n", now)

		var activeBookings int64
		if err = db.Table("booking_to_orders").
			Where("employee_id = ? AND order_end > ?", workerID, now).
			Count(&activeBookings).Error; err != nil {
			fmt.Printf("Error checking active bookings: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error checking active bookings.",
			})
		}

		fmt.Printf("Active bookings for employee %d: %d\n", workerID, activeBookings)

		if activeBookings > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "У данного сотрудника есть предстоящие заказы.",
			})
		}

		if err = db.Model(&employee).Update("status", "Уволен").Error; err != nil {
			fmt.Printf("Error updating employee status: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error updating employee status.",
			})
		}

		fmt.Printf("Employee %d fired successfully.\n", workerID)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Employee has been fired successfully.",
			"worker": fiber.Map{
				"employee_id":   employee.EmployeeID,
				"full_name":     employee.FullName,
				"position":      employee.Position,
				"hire_date":     employee.HireDate,
				"birth_date":    employee.BirthDate,
				"passport_data": employee.PassportData,
				"phone_number":  employee.PhoneNumber,
				"status":        employee.Status,
			},
		})
	}
}
