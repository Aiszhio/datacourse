package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type CustomPurchase struct {
	PurchaseID int           `gorm:"column:purchase_id;primaryKey;autoIncrement" json:"purchase_id"`
	MaterialID int           `gorm:"column:material_id;not null" json:"material_id"`
	Cost       CustomFloat64 `gorm:"column:cost;not null" json:"cost"`
	Supplier   string        `gorm:"column:supplier;not null" json:"supplier"`
	Quantity   int           `gorm:"column:quantity;not null" json:"quantity"`
	SupplyDate CustomTime    `gorm:"column:supply_date;not null" json:"supply_date"`
}

func AddPurchase(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var purchaseItem CustomPurchase

		if err := c.BodyParser(&purchaseItem); err != nil {
			fmt.Println("error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка получения данных",
			})
		}

		loc, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			fmt.Println("error loading location:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка загрузки временной зоны",
			})
		}

		parsedTime := purchaseItem.SupplyDate.Time.In(loc)

		purchaseItem.SupplyDate = CustomTime{parsedTime}

		tx := db.Begin()
		if tx.Error != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при начале транзакции",
			})
		}

		if err = tx.Table("material_purchases").Create(&purchaseItem).Error; err != nil {
			fmt.Println("error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Не удалось совершить покупку",
			})
		}

		var purchase db2.Material
		if err = tx.Table("materials").Where("material_id = ?", purchaseItem.MaterialID).First(&purchase).Error; err != nil {
			tx.Rollback()
			fmt.Println("error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка парсинга материалов",
			})
		}

		newQuantity := purchase.Quantity + purchaseItem.Quantity
		if err = tx.Model(&purchase).Update("quantity", newQuantity).Error; err != nil {
			tx.Rollback()
			fmt.Println("error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка обновления количества материалов",
			})
		}

		if err = tx.Commit().Error; err != nil {
			fmt.Println("error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при сохранении изменений в базе данных",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Закупка успешно проведена!",
		})
	}
}
