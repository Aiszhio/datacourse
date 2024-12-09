package handlers

import (
	"fmt"
	db2 "github.com/Aiszhio/datacourse.git/pkg/db"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

func AddClient(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var client db2.Client

		if err := c.BodyParser(&client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Сервер не смог получить данные",
			})
		}

		if err := ValidateClient(client); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Table("clients").Create(&client).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Пользователь не был записан в БД",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Пользователь был успешно добавлен",
		})
	}
}

func GetClients(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var clientList []db2.Client

		err := db.Table("clients").Find(&clientList).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Не нашли данные в таблице клиентов",
			})
		}

		if len(clientList) == 0 {
			fmt.Println("Ничего нет")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"clients": clientList,
		})
	}
}

func ValidateClient(client db2.Client) error {
	namePattern := `^[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?(?:\s[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?)?(?:\s[А-ЯЁ][а-яё]+(?:-[А-ЯЁ][а-яё]+)?)?$`
	phonePattern := `^[7][0-9]{10}$`
	emailPattern := `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`
	reName := regexp.MustCompile(namePattern)
	rePhone := regexp.MustCompile(phonePattern)
	reEmail := regexp.MustCompile(emailPattern)

	client.FullName = strings.TrimSpace(client.FullName)
	client.PhoneNumber = strings.TrimSpace(client.PhoneNumber)
	client.Email = strings.TrimSpace(client.Email)

	if !reName.MatchString(client.FullName) {
		fmt.Println(client.FullName)
		return fmt.Errorf("Имя пользователя неверно.")
	}

	if !rePhone.MatchString(client.PhoneNumber) {
		fmt.Println(client.PhoneNumber)
		return fmt.Errorf("Номер телефона пользователя неправильный.")
	}

	if !reEmail.MatchString(client.Email) {
		fmt.Println(client.Email)
		return fmt.Errorf("Почта пользователя неверна.")
	}

	if len(client.PhoneNumber) != 11 {
		return fmt.Errorf("Номер телефона должен содержать ровно 11 символов.")
	} else if len(client.FullName) > 80 {
		return fmt.Errorf("Имя клиента слишком длинное.")
	}

	if !reEmail.MatchString(client.Email) {
		return fmt.Errorf("Электронная почта указана неверно.")
	}

	return nil
}
