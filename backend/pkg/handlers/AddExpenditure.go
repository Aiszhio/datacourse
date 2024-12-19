package handlers

import (
	"github.com/Aiszhio/datacourse.git/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type CustomExpenditure struct {
	ExpenditureID   int       `gorm:"column:expenditure_id;primaryKey;autoIncrement" json:"expenditure_id"`
	MaterialID      int       `gorm:"column:material_id;not null" json:"material_id"`
	ExpenditureDate time.Time `gorm:"column:expenditure_date;not null" json:"expenditure_date"`
	Quantity        int       `gorm:"column:quantity;not null" json:"quantity"`
}

func AddExpenditure(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var expenditure CustomExpenditure

		if err := c.BodyParser(&expenditure); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга данных",
			})
		}

		if err := utils.CheckWorkingHours(expenditure.ExpenditureDate); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		tx := db.Begin()

		if tx.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка транзакции",
			})
		}

		if err := tx.Table("material_expenditures").Create(&expenditure).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка создания траты",
			})
		}

		var existingQuantity int
		if err := tx.Table("materials").Select("quantity").Where("material_id = ?", expenditure.MaterialID).Scan(&existingQuantity).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при просмотре таблицы материалов",
			})
		}

		if existingQuantity < expenditure.Quantity {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Недостаточно материала для расхода",
			})
		}

		newQuantity := existingQuantity - expenditure.Quantity
		if err := tx.Table("materials").Where("material_id = ?", expenditure.MaterialID).Update("quantity", newQuantity).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при обновлении таблицы материалов",
			})
		}

		if tx.Commit().Error != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при обновлении данных",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Вы успешно провели операцию траты!",
		})
	}
}
