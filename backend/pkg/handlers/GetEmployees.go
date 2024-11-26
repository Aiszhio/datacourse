package handlers

import (
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strings"
)

func GetEmployees(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var employeesList []database.Employee

		if err := db.Table("employees").Find(&employeesList).Error; err != nil {
			return err
		}

		for i := range employeesList {
			passportData := strings.TrimSpace(employeesList[i].PassportData)
			employeesList[i].PassportData = passportData[:4] + " " + passportData[4:]
		}

		return c.JSON(fiber.Map{
			"employees": employeesList,
		})
	}
}
