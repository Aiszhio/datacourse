package main

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Authentification"
	"github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/Aiszhio/datacourse.git/pkg/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	webApp.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "*",
	}))
	webApp.Post("/api/login", Authentification.Authorize(dbu))
	webApp.Get("/client", handlers.GetUserData(dbu))

	defer log.Fatal(webApp.Listen(":8080"))
}
