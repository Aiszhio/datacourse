package handlers

import (
	"errors"
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"regexp"
	"strconv"
)

type ServiceToUpdate struct {
	Price float64 `gorm:"column:price;not null" json:"price"`
	Name  string  `gorm:"column:name" json:"name"`
}

func UpdateService(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var updatedService ServiceToUpdate
		serviceIDStr := c.Params("service_id")
		if serviceIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Необходим идентификатор",
			})
		}

		serviceID, err := strconv.Atoi(serviceIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка конвертации данных",
			})
		}

		if err = c.BodyParser(&updatedService); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка обработки данных",
			})
		}

		var existingService db2.Service
		if err = db.Table("services").First(&existingService, "service_id = ?", serviceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println(err)
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Данной услуги не обнаружено",
				})
			}
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка",
			})
		}

		if err = ValidatePutService(updatedService); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err = db.Model(&existingService).Updates(updatedService).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка сохранения данных в БД",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Данные о услуге успешно обновлены!",
			"service": updatedService,
		})
	}
}

func ValidatePutService(service ServiceToUpdate) error {
	validNames := []string{"Печать", "Фотокнига", "Фотоальбом", "Фотосессия", "Съемка", "Аренда"}

	for _, name := range validNames {
		matched, _ := regexp.MatchString("(?i)"+name, service.Name)
		if matched {
			return nil
		}
	}

	return errors.New("Название услуги должно содержать одно из следующих слов: Печать, Фотокнига, Фотоальбом, Фотосессия, Съемка")
}
