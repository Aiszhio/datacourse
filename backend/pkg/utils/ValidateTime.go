package utils

import (
	"fmt"
	"time"
)

func CheckWorkingHours(t time.Time) error {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return fmt.Errorf("невозможно загрузить временную зону Москвы: %v", err)
	}

	moscowTime := t.In(loc)
	t1 := moscowTime.Add(time.Hour)

	start := time.Date(moscowTime.Year(), moscowTime.Month(), moscowTime.Day(), 9, 0, 0, 0, loc)
	end := time.Date(moscowTime.Year(), moscowTime.Month(), moscowTime.Day(), 19, 0, 0, 0, loc)

	if moscowTime.Before(start) || t1.After(end) {
		return fmt.Errorf("интервал [%s - %s] выходит за пределы рабочего времени (09:00 - 19:00 МСК)",
			moscowTime.Format("15:04"), t1.Format("15:04"))
	}

	return nil
}
