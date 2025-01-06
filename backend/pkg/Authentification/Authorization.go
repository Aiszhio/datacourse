package Authentification

import (
	"errors"
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/Aiszhio/datacourse.git/pkg/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Customer struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func getRoleByPhone(phone, password string, db *gorm.DB, model interface{}, role, idColumn string) (string, error) {
	if err := db.Where("phone_number = ? AND "+idColumn+" = ?", phone, password).First(model).Error; err != nil {
		return "", errors.New("not found")
	}

	if role == "client" {
		return "client", nil
	}

	var position string
	err := db.Table("employees").Select("position").Where("phone_number = ? AND position IN ?", phone, []string{"Фотограф", "Администратор"}).Scan(&position).Error
	if err != nil {
		return "", err
	}

	switch position {
	case "Фотограф":
		return "worker", nil
	case "Администратор":
		return "admin", nil
	default:
		return "", errors.New("unknown role")
	}
}

func userDefineRole(phone, password string, db *gorm.DB) (string, error) {
	if role, err := getRoleByPhone(phone, password, db, &database.Client{}, "client", "client_id"); err == nil {
		return role, nil
	}

	if role, err := getRoleByPhone(phone, password, db, &database.Employee{}, "worker", "employee_id"); err == nil {
		return role, nil
	}

	if role, err := getRoleByPhone(phone, password, db, &database.Employee{}, "admin", "employee_id"); err == nil {
		return role, nil
	}

	return "", fmt.Errorf("Пользователь с таким номером телефона или паролем не найден")
}

func Authorize(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var customer Customer
		var userInfo interface{}

		if err := c.BodyParser(&customer); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request format",
			})
		}

		role, err := userDefineRole(customer.Phone, customer.Password, db)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = Redis.SetKey(client, "phone", customer.Phone)
		fmt.Println(customer.Phone)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error saving phone to Redis",
			})
		}

		userInfo, err = handlers.GetUserByPhone(client, db, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error retrieving user info",
			})
		}

		err = Redis.SetKey(client, "role", role)
		fmt.Println(Redis.GetKey(client, "role"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error saving role to Redis",
			})
		}

		err = Redis.SetMultipleKey(client, "UserInfo", userInfo)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error saving user info to Redis",
			})
		}

		err = SetCookie(client, c)
		fmt.Println(Redis.GetKey(client, "sessionID"))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error saving cookies to Redis",
			})
		}

		var employee database.Employee

		if role != "client" {
			if err = db.Table("employees").Where("phone_number = ?", customer.Phone).
				First(&employee).Error; err != nil {
				fmt.Println("2312312312312312312")
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Ошибка при логине",
				})
			}
			if employee.Status == "Уволен" {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Вы уволены",
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"role":    role,
			"message": "Login successful",
		})
	}
}
