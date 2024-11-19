package handlers

import (
	"fmt"
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetEmployees(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var employeesList []database.Employee

		if err := db.Table("employees").Find(&employeesList).Error; err != nil {
			return err
		}

		fmt.Println(employeesList)

		return c.JSON(fiber.Map{
			"employees": employeesList,
		})
	}
}
