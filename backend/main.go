package main

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Authentification"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
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

	rdb := Redis.InitRedisClient()

	if err = db.MigrateDB(); err != nil {
		fmt.Println("Error during DB migration:", err)
		return
	}

	if err = db.InsertUniqueDB(dbu); err != nil {
		fmt.Println("Error inserting unique data:", err)
		return
	}

	webApp.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowHeaders:     "Content-Type, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	webApp.Use(Authentification.SessionMiddleware(rdb))

	webApp.Post("/api/login", Authentification.Authorize(dbu, rdb))
	webApp.Get("/api/CheckRole", Authentification.CheckUserRole(rdb))
	webApp.Get("/api/user", handlers.GetUserData(rdb))
	webApp.Get("/api/orders", handlers.ClientOrdersApi(dbu, rdb))
	webApp.Get("/api/services", handlers.GetServices(dbu))
	webApp.Post("/api/createOrder", handlers.CreateOrder(dbu, rdb))
	webApp.Get("/api/employees", handlers.GetEmployees(dbu))
	webApp.Post("/api/employees", handlers.AddEmployee(dbu))

	defer rdb.Close()
	defer log.Fatal(webApp.Listen(":8080"))
}
