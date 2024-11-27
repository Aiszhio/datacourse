package handlers

import (
	"errors"
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

func AddEmployee(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var employee database.Employee

		fmt.Println("Получение данных для добавления сотрудника...")

		if err := c.BodyParser(&employee); err != nil {
			fmt.Println("Ошибка при парсинге тела запроса:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Проверьте данные.",
			})
		}

		fmt.Printf("Полученные данные сотрудника: %+v\n", employee)

		if strings.TrimSpace(employee.Status) == "" {
			employee.Status = "Работает"
		}

		if employee.FullName == "" || employee.Position == "" || employee.HireDate.IsZero() || employee.BirthDate.IsZero() || employee.PassportData == "" || employee.PhoneNumber == "" {
			fmt.Println("Некоторые поля отсутствуют:", employee)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Необходимо ввести все поля.",
			})
		}

		if err := ValidateEmployeeData(employee); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		var passportCheck database.Employee
		if err := db.Where("passport_data = ?", employee.PassportData).First(&passportCheck).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Проверьте паспортные данные.",
			})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка проверки паспортных данных.",
			})
		}

		var phoneCheck database.Employee
		if err := db.Where("phone_number = ?", employee.PhoneNumber).First(&phoneCheck).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Сотрудник с таким номером телефона уже существует.",
			})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка проверки номера телефона.",
			})
		}

		if err := db.Create(&employee).Error; err != nil {
			fmt.Println("Ошибка внесения сотрудника в БД:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Ошибка внесения сотрудника в БД.",
			})
		}

		formattedHireDate := employee.HireDate.Format("2006-01-02")
		formattedBirthDate := employee.BirthDate.Format("2006-01-02")

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Сотрудник успешно добавлен.",
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

func ValidateEmployeeData(employee database.Employee) error {
	namePattern := `^[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?(?:\s[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?)?(?:\s[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?)?$`
	phonePattern := `^[7][0-9]{10}$`
	passportPattern := `^[0-9]{10}$`
	reName := regexp.MustCompile(namePattern)
	rePhone := regexp.MustCompile(phonePattern)
	rePassport := regexp.MustCompile(passportPattern)

	employee.FullName = strings.TrimSpace(employee.FullName)
	employee.PassportData = strings.TrimSpace(employee.PassportData)
	employee.PhoneNumber = strings.TrimSpace(employee.PhoneNumber)

	if !reName.MatchString(employee.FullName) {
		return fmt.Errorf("Имя сотрудника неверно.")
	}

	if !rePassport.MatchString(employee.PassportData) {
		return fmt.Errorf("Проверьте паспортные данные.")
	}

	if !rePhone.MatchString(employee.PhoneNumber) {
		return fmt.Errorf("Номер телефона сотрудника неправильный.")
	}

	if len(employee.PassportData) != 10 {
		return fmt.Errorf("Паспортные данные должны содержать ровно 10 символов.")
	} else if len(employee.PhoneNumber) != 11 {
		return fmt.Errorf("Номер телефона должен содержать ровно 11 символов.")
	} else if len(employee.FullName) > 80 {
		return fmt.Errorf("Имя сотрудника слишком длинное.")
	}

	today := time.Now()
	years := today.Year() - employee.BirthDate.Year()

	if today.Month() < employee.BirthDate.Month() || (today.Month() == employee.BirthDate.Month() && today.Day() < employee.BirthDate.Day()) {
		years--
	}

	if years < 18 {
		return fmt.Errorf("Сотрудник должен быть старше 18 лет.")
	}

	validStatuses := map[string]bool{
		"Работает": true,
		"Уволен":   true,
	}

	if employee.Status != "" {
		if !validStatuses[employee.Status] {
			return fmt.Errorf("Неверный статус сотрудника.")
		}
	}

	return nil
}
