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
		fmt.Println("fsdfsd")
		if err := c.BodyParser(&employee); err != nil {
			fmt.Println(err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body.",
			})
		}

		if employee.FullName == "" || employee.Position == "" || employee.HireDate.IsZero() || employee.BirthDate.IsZero() || employee.PassportData == "" || employee.PhoneNumber == "" {
			fmt.Println(employee)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "All fields are required.",
			})
		}

		if err := ValidateEmployeeData(employee, c); err != nil {
			return err
		}

		var passportCheck database.Employee
		if err := db.Where("passport_data = ?", employee.PassportData).First(&passportCheck).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Employee with this passport already exists.",
			})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error checking passport data.",
			})
		}

		var phoneCheck database.Employee
		if err := db.Where("phone_number = ?", employee.PhoneNumber).First(&phoneCheck).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Employee with this phone number already exists.",
			})
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error checking phone number.",
			})
		}

		if err := db.Create(&employee).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error in creating employee.",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Employee added successfully.",
		})
	}
}

func ValidateEmployeeData(employee database.Employee, c *fiber.Ctx) error {
	namePattern := `^[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?(?:\s[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?)?(?:\s[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?)?$`
	phonePattern := `^[7][0-9]{10}$`
	passportPattern := `^[0-9]{10}$`
	reName := regexp.MustCompile(namePattern)
	rePhone := regexp.MustCompile(phonePattern)
	rePassport := regexp.MustCompile(passportPattern)

	employee.FullName = strings.TrimSpace(employee.FullName)
	if !reName.MatchString(employee.FullName) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Employee name isn't correct. ",
		})
	}

	employee.PassportData = strings.TrimSpace(employee.PassportData)
	if !rePassport.MatchString(employee.PassportData) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Employee passport data isn't correct. ",
		})
	}

	employee.PhoneNumber = strings.TrimSpace(employee.PhoneNumber)
	if !rePhone.MatchString(employee.PhoneNumber) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Employee phone number isn't correct. ",
		})
	}

	today := time.Now()
	years := today.Year() - employee.BirthDate.Year()

	if today.Month() < employee.BirthDate.Month() || (today.Month() == employee.BirthDate.Month() && today.Day() < employee.BirthDate.Day()) {
		years--
	}

	if years < 18 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Employee Birthdate is wrong",
		})
	}

	return nil
}
