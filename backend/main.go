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
	webApp.Post("/api/register", Authentification.Registration(dbu))
	webApp.Get("/api/CheckRole", Authentification.CheckUserRole(dbu, rdb))
	webApp.Get("/api/user", handlers.GetUserName(rdb))
	webApp.Get("/api/orders", handlers.ClientOrdersApi(dbu, rdb))
	webApp.Get("/api/bookings", handlers.GetBookings(dbu, rdb))
	webApp.Get("/api/orders/worker", handlers.GetWorkerOrders(dbu, rdb))
	webApp.Get("/api/services", handlers.GetServices(dbu))
	webApp.Post("/api/createOrder", handlers.CreateBooking(dbu, rdb))
	webApp.Get("/api/employees", handlers.GetEmployees(dbu))
	webApp.Post("/api/employees", handlers.AddEmployee(dbu))
	webApp.Get("/api/orders/admin", handlers.GetAdminOrders(dbu))
	webApp.Post("/api/orders/admin", handlers.AddOrder(dbu))
	webApp.Get("/api/materials", handlers.GetMaterials(dbu))
	webApp.Get("/api/expenditures", handlers.GetMaterialExpenditures(dbu))
	webApp.Post("/api/expenditures", handlers.AddExpenditure(dbu))
	webApp.Post("/api/purchases", handlers.AddPurchase(dbu))
	webApp.Get("/api/purchases", handlers.GetMaterialPurchases(dbu))
	webApp.Get("/api/services/admin", handlers.GetServicesAndEquipment(dbu))
	webApp.Post("/api/services/admin", handlers.AddService(dbu))
	webApp.Post("/api/equipment/admin", handlers.AddEquipment(dbu))
	webApp.Get("/api/bookings/admin", handlers.GetAdminBookings(dbu))
	webApp.Delete("/api/bookings/:bookingID", handlers.DeleteBooking(dbu))
	webApp.Put("/api/employees/:id/fire", handlers.FireWorker(dbu))
	webApp.Put("/api/employees/:id", handlers.ChangeEmployee(dbu))
	webApp.Post("/api/createOrder/admin", handlers.AdminCreateBooking(dbu))
	webApp.Delete("/api/bookings/admin/:bookingId", handlers.DeleteAdminBooking(dbu))
	webApp.Put("/api/services/:service_id", handlers.UpdateService(dbu))
	webApp.Put("/api/equipment/:equipment_id", handlers.UpdateEquipment(dbu))
	webApp.Get("/api/admin/clients", handlers.GetClients(dbu))
	webApp.Post("/api/admin/clients", handlers.AddClient(dbu))
	webApp.Delete("/api/equipment/:equipmentID", handlers.DeleteEquipment(dbu))
	webApp.Delete("/api/services/:service_id", handlers.DeleteService(dbu))
	webApp.Put("/api/admin/clients/:id", handlers.PutClient(dbu))
	webApp.Delete("/api/admin/clients/:id", handlers.DeleteClient(dbu))
	webApp.Delete("/api/orders/:orderId", handlers.DeleteOrder(dbu))
	webApp.Post("/api/orders/issue/:id", handlers.IssueOrder(dbu))

	defer rdb.Close()
	defer log.Fatal(webApp.Listen(":8080"))
}
