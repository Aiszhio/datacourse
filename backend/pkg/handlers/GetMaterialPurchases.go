package handlers

import (
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type PurchasesWithNames struct {
	PurchaseID   int       `json:"purchase_id"`
	MaterialName string    `json:"material_name"`
	Cost         float64   `json:"cost"`
	Supplier     string    `json:"supplier"`
	Quantity     int       `json:"quantity"`
	SupplyDate   time.Time `json:"supply_date"`
}

func GetMaterialPurchases(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var purchasesList []database.MaterialPurchase

		err := db.Table("material_purchases").Find(&purchasesList).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		newPurchasesList := GetPurchasesNames(purchasesList, db)

		if len(newPurchasesList) == 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": "Error in searching the material purchases",
			})
		}

		return c.JSON(fiber.Map{
			"purchases": newPurchasesList,
		})
	}
}

func GetPurchasesNames(purchasesList []database.MaterialPurchase, db *gorm.DB) []PurchasesWithNames {
	var newPurchasesList []PurchasesWithNames

	for _, purchase := range purchasesList {
		var tempMaterialName string

		err := db.Table("materials").Select("material_name").
			Where("material_id = ?", purchase.MaterialID).Scan(&tempMaterialName).Error
		if err != nil {
			return nil
		}

		newPurchasesList = append(newPurchasesList, PurchasesWithNames{
			PurchaseID:   purchase.PurchaseID,
			MaterialName: tempMaterialName,
			Cost:         purchase.Cost,
			Supplier:     purchase.Supplier,
			Quantity:     purchase.Quantity,
			SupplyDate:   purchase.SupplyDate,
		})
	}

	return newPurchasesList
}
