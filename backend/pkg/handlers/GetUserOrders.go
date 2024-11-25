package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetUserByPhone(rdb *redis.Client, db *gorm.DB, role string) (interface{}, error) {
	userInfo := make(map[string]interface{})
	var tableName string
	if role == "client" {
		tableName = "clients"
	} else {
		tableName = "employees"
	}

	phone, err := Redis.GetKey(rdb, "phone")
	if err != nil {
		return nil, fmt.Errorf("error in getting phone number from Redis: %w", err)
	}

	err = db.Table(tableName).Where("phone_number = ?", phone).Limit(1).Find(&userInfo).Error
	if err != nil {
		return nil, fmt.Errorf("error in getting user info from table %s: %w", tableName, err)
	}

	return userInfo, nil
}

func GetUserName(client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userInfo, err := Redis.GetMultipleKey(client, "UserInfo")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "cannot get user",
			})
		}

		name, ok := userInfo["full_name"].(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Full Name is not found",
			})
		}

		return c.JSON(fiber.Map{
			"name": name,
		})
	}
}
