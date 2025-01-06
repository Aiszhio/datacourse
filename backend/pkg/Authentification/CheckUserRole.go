package Authentification

import (
	"errors"
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"strconv"
)

func CheckUserRole(db *gorm.DB, client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := Redis.GetKey(client, "role")

		if err != nil {
			return err
		}

		userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if role != "client" {

			var userID int

			switch v := userInfo["employee_id"].(type) {
			case int:
				userID = v
			case float64:
				userID = int(v)
			case string:
				userID, err = strconv.Atoi(v)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "invalid type of user id",
				})
			}

			var employee db2.Employee
			if err = db.Table("employees").Where("employee_id = ? AND status = ?", userID, "Уволен").First(&employee).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
				} else {
					fmt.Println("Ошибка при проверке статуса пользователя:", err)
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": "Внутренняя ошибка сервера",
					})
				}
			} else {
				fmt.Println("Пользователь уволен")
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Доступ запрещен, вы уволены",
				})
			}
		}

		return c.JSON(fiber.Map{
			"role": role,
		})
	}
}
