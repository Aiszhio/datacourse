package Authentification

import (
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func CheckUserRole(client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, err := Redis.GetKey(client, "role")
		if err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"role": role,
		})
	}
}
