package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func DeleteService(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		serviceIDStr := c.Params("service_id")
		serviceID, err := strconv.Atoi(serviceIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Некорректный ID услуги",
			})
		}

		if err = ValidateDeleteService(serviceID, db); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tx := db.Begin()
		if tx.Error != nil {
			fmt.Println("Ошибка начала транзакции:", tx.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при начале транзакции",
			})
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				fmt.Println("Паника при транзакции:", r)
			}
		}()

		result := tx.Delete(&db2.Service{}, serviceID)
		if result.Error != nil {
			tx.Rollback()
			fmt.Println("Ошибка при удалении услуги:", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Не удалось удалить услугу",
			})
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Услуга не найдена",
			})
		}

		if err = tx.Commit().Error; err != nil {
			fmt.Println("Ошибка коммита транзакции:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при завершении транзакции",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Услуга успешно удалена",
		})
	}
}
