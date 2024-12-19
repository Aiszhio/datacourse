package handlers

import (
	"errors"
	"fmt"
	"strconv"

	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UpdateClientInput struct {
	FullName    string `json:"FullName,omitempty"`
	PhoneNumber string `json:"PhoneNumber,omitempty"`
	Email       string `json:"Email,omitempty"`
}

func PutClient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		if idStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Требуется идентификатор клиента",
			})
		}

		fmt.Println("Вот айдишник ", idStr)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка конвертации id клиента",
			})
		}

		var existingClient db2.Client
		if err = db.First(&existingClient, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Клиент не найден",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при поиске клиента",
			})
		}

		var updateData UpdateClientInput
		if err = c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка в парсинге данных",
			})
		}

		if updateData.FullName != "" {
			existingClient.FullName = updateData.FullName
		}
		if updateData.PhoneNumber != "" {
			existingClient.PhoneNumber = updateData.PhoneNumber
		}
		if updateData.Email != "" {
			existingClient.Email = updateData.Email
		}

		if err = ValidateClient(existingClient); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err = db.Save(&existingClient).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Клиент с таким номером телефона или почтой уже существует",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Ошибка при обновлении клиента",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Клиент успешно обновлен",
		})
	}
}
