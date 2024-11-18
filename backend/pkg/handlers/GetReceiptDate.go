package handlers

import (
	"fmt"
	"time"
)

func GetReceiptDate(serviceName string, orderDate time.Time) (time.Time, error) {
	switch serviceName {
	case "Аренда студии (фотосессия)":
		return orderDate.Add(time.Hour * 1), nil
	case "Печать фотографий":
		return orderDate.Add(time.Hour * 24), nil
	case "Создание фотокниги":
		return orderDate.AddDate(0, 0, 7), nil
	case "Создание фотоальбома":
		return orderDate.AddDate(0, 0, 5), nil
	case "Печать на холсте":
		return orderDate.AddDate(0, 0, 3), nil
	case "Цифровая обработка фотографий":
		return orderDate.Add(time.Hour * 48), nil
	case "Профессиональная портретная съемка":
		return orderDate.Add(time.Hour * 4), nil
	case "Свадебная фотосессия":
		return orderDate.AddDate(0, 0, 14), nil
	case "Фотографирование детских праздников":
		return orderDate.Add(time.Hour * 6), nil
	case "Съемка мероприятий":
		return orderDate.Add(time.Hour * 8), nil
	default:
		return time.Time{}, fmt.Errorf("unknown service name: %s", serviceName)
	}
}
