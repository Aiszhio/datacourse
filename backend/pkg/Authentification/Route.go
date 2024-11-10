package Authentification

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func DashBoardHandler(c *fiber.Ctx) error {
	role := c.Locals("role").(string)

	switch role {
	case "admin":
		return c.Redirect("/admin/")
	case "worker":
		return c.Redirect("/worker/")
	case "client":
		return c.Redirect("/client/")
	}

	return c.SendString(fmt.Sprintf("Welcome to the DashBoard, Role: %s", role))
}
