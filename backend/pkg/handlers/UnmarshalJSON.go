package handlers

import (
	"encoding/json"
	"fmt"
	"time"
)

func (u *UserBooking) UnmarshalJSON(data []byte) error {
	type Alias UserBooking
	aux := &struct {
		BookingTime string `json:"booking_time"`
		*Alias
	}{
		Alias: (*Alias)(u),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("invalid JSON structure: %v", err)
	}

	parsedTime, err := time.Parse(time.RFC3339, aux.BookingTime)
	if err != nil {
		return fmt.Errorf("invalid date format for booking_time: %v", err)
	}

	u.BookingTime = parsedTime
	return nil
}
