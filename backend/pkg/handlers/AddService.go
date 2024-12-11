package handlers

import (
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ServiceInput struct {
	Price             float64 `json:"price"`
	Name              string  `json:"name"`
	RequiredEquipment []int   `json:"requiredEquipment"`
}

func AddService(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input ServiceInput

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка при парсинге данных",
			})
		}

		if err := ValidateService(input, db); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		serviceDB := db2.Service{
			Price: input.Price,
			Name:  input.Name,
		}

		tx := db.Begin()
		if tx.Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка в начале транзакции",
			})
		}

		if err := tx.Create(&serviceDB).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при записи данных в таблицу услуг",
			})
		}

		if len(input.RequiredEquipment) > 0 {
			var serviceRequirements []db2.ServiceRequirement
			for _, eqID := range input.RequiredEquipment {
				serviceRequirement := db2.ServiceRequirement{
					ServiceID:   serviceDB.ServiceID,
					EquipmentID: eqID,
				}
				serviceRequirements = append(serviceRequirements, serviceRequirement)
			}

			if err := tx.Create(&serviceRequirements).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Ошибка при записи требований к услуге",
				})
			}
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при завершении транзакции",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Услуга успешно добавлена",
			"service": serviceDB,
		})
	}
}
