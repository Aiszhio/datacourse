package Authentification

import (
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

func SetCookie(client *redis.Client, c *fiber.Ctx) error {
	sessionID := uuid.New().String()

	err := Redis.SetKey(client, "sessionID", sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error saving session to Redis")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "sessionID",
		Value:    sessionID,
		Expires:  time.Now().Add(30 * time.Minute),
		HTTPOnly: true,
		Path:     "/",
		Domain:   "localhost",
		SameSite: "None",
	})

	return nil
}
