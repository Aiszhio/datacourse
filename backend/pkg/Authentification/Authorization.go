package Authentification

import (
	"errors"
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Customer struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func getRoleByPhone(phone, password string, db *gorm.DB, model interface{}, role, idColumn string) (string, error) {
	if err := db.Where("phone_number = ? AND "+idColumn+" = ?", phone, password).First(model).Error; err == nil {
		return role, nil
	}
	return "", errors.New("not found")
}

func userDefine(phone, password string, db *gorm.DB) (string, error) {
	if role, err := getRoleByPhone(phone, password, db, &database.Client{}, "client", "client_id"); err == nil {
		return role, nil
	}

	if role, err := getRoleByPhone(phone, password, db, &database.Employee{}, "worker", "employee_id"); err == nil {
		return role, nil
	}

	if role, err := getRoleByPhone(phone, password, db, &database.Employee{}, "admin", "employee_id"); err == nil {
		return role, nil
	}

	return "", fmt.Errorf("user with phone %s not found in any role", phone)
}

func Authorize(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var customer Customer

		if err := c.BodyParser(&customer); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		role, err := userDefine(customer.Phone, customer.Password, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"role": role})
	}
}
