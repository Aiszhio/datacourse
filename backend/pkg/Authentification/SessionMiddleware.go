package Authentification

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func SessionMiddleware(client *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/" || c.Path() == "/api/login" {
			return c.Next()
		}

		sessionID := c.Cookies("sessionID")
		if sessionID == "" {
			fmt.Println("Session ID not found in cookies")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: No session ID",
			})
		}

		role, err := Redis.GetKey(client, "sessionID")
		if err != nil || role == "" {
			fmt.Println("Error retrieving role from Redis for session:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Session expired or invalid",
			})
		}

		c.Locals("role", role)
		fmt.Println("User role retrieved from Redis:", role)

		return c.Next()
	}
}
