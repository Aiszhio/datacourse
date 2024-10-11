package main

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	webApp := fiber.New()

	dbu, err := db.InitDB()
	if err != nil {
		fmt.Println("Error initializing DB:", err)
		return
	}

	if err = db.MigrateDB(); err != nil {
		fmt.Println("Error during DB migration:", err)
		return
	}

	if err = db.InsertUniqueDB(dbu); err != nil {
		fmt.Println("Error inserting unique data:", err)
		return
	}

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"key": "here we are!!!"})
	})
	fmt.Println("Listening on port 8080")
	webApp.Get("/data", func(c *fiber.Ctx) error {
		data, err := db.GetAllData(dbu)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(data)
	})

	defer log.Fatal(webApp.Listen(":8080"))
}
