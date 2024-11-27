package handlers

import (
	"errors"
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

func ChangeEmployee(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		workerIDStr := c.Params("id")
		if workerIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Необходим идентификатор сотрудника.",
			})
		}

		workerID, err := strconv.Atoi(workerIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Неправильный идентификатор сотрудника.",
			})
		}

		var updatedData struct {
			FullName     string `json:"full_name"`
			Position     string `json:"position"`
			HireDate     string `json:"hire_date"`
			BirthDate    string `json:"birth_date"`
			PassportData string `json:"passport_data"`
			PhoneNumber  string `json:"phone_number"`
			Status       string `json:"status"`
		}

		if err = c.BodyParser(&updatedData); err != nil {
			fmt.Println("Ошибка при парсинге тела запроса:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Ошибка обработки данных.",
			})
		}

		var employee database.Employee
		if err = db.First(&employee, workerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"message": "Сотрудник не найден.",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка обновления сотрудника.",
			})
		}

		if updatedData.FullName != "" {
			employee.FullName = updatedData.FullName
		}
		if updatedData.Position != "" {
			employee.Position = updatedData.Position
		}

		var hireDate, birthDate time.Time

		if updatedData.HireDate != "" {
			hireDate, err = time.Parse(time.RFC3339, updatedData.HireDate)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Неправильный формат даты оформления.",
				})
			}
			employee.HireDate = hireDate
		}

		if updatedData.BirthDate != "" {
			birthDate, err = time.Parse(time.RFC3339, updatedData.BirthDate)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Неправильный формат даты рождения.",
				})
			}
			employee.BirthDate = birthDate
		}

		if updatedData.PassportData != "" {
			employee.PassportData = strings.ReplaceAll(updatedData.PassportData, " ", "")
		}
		if updatedData.PhoneNumber != "" {
			employee.PhoneNumber = strings.ReplaceAll(updatedData.PhoneNumber, "-", "")
		}
		if updatedData.Status != "" {
			employee.Status = updatedData.Status
		}

		if err = ValidateEmployeeData(employee); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		if err = db.Save(&employee).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка сохранения данных сотрудника.",
			})
		}

		formattedHireDate := employee.HireDate.Format("2006-01-02")
		formattedBirthDate := employee.BirthDate.Format("2006-01-02")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Данные сотрудника успешно изменены!",
			"worker": fiber.Map{
				"id":            employee.EmployeeID,
				"full_name":     employee.FullName,
				"position":      employee.Position,
				"hire_date":     formattedHireDate,
				"birth_date":    formattedBirthDate,
				"passport_data": employee.PassportData,
				"phone_number":  employee.PhoneNumber,
				"status":        employee.Status,
			},
		})
	}
}
