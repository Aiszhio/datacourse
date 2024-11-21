package handlers

import (
	database "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type ExpendituresWithNames struct {
	ExpenditureID   int       `json:"expenditure_id"`
	MaterialName    string    `json:"material_name"`
	ExpenditureDate time.Time `json:"expenditure_date"`
	Quantity        int       `json:"quantity"`
}

func GetMaterialExpenditures(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var expendituresList []database.MaterialExpenditure

		err := db.Table("material_expenditures").Find(&expendituresList).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		newExpendituresList := GetExpendituresNames(expendituresList, db)

		if len(newExpendituresList) == 0 {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"error": "Error in searching the material expenditures",
			})
		}

		return c.JSON(fiber.Map{
			"expenditures": newExpendituresList,
		})
	}
}

func GetExpendituresNames(expendituresList []database.MaterialExpenditure, db *gorm.DB) []ExpendituresWithNames {
	var newExpendituresList []ExpendituresWithNames

	for i := 0; i < len(expendituresList); i++ {
		var tempExpenditureName string

		err := db.Table("materials").Select("material_name").
			Where("material_id = ?", expendituresList[i].MaterialID).Scan(&tempExpenditureName).Error
		if err != nil {
			return nil
		}

		newExpendituresList = append(newExpendituresList, ExpendituresWithNames{
			ExpenditureID:   expendituresList[i].ExpenditureID,
			MaterialName:    tempExpenditureName,
			ExpenditureDate: expendituresList[i].ExpenditureDate,
			Quantity:        expendituresList[i].Quantity,
		})
	}

	return newExpendituresList
}
