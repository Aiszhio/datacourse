package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

func DeleteEquipment(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("equipmentID")

		idInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Ошибка парсинга ID оборудования:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Некорректный ID оборудования",
			})
		}

		if err = ValidateDeleteEquipment(idInt, db); err != nil {
			fmt.Println("Валидация оборудования не пройдена:", err)
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

		if err = tx.Table("equipment").Where("equipment_id = ?", idInt).Delete(nil).Error; err != nil {
			tx.Rollback()
			fmt.Println("Ошибка при удалении оборудования:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Не удалось удалить оборудование",
			})
		}

		if err = tx.Commit().Error; err != nil {
			tx.Rollback()
			fmt.Println("Ошибка коммита транзакции:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при завершении транзакции",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Оборудование успешно списано",
		})
	}
}
