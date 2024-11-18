package handlers

import (
	"encoding/json"
	"fmt"
	"time"
)

func (order *ClientOrder) UnmarshalJSON(data []byte) error {
	type Alias ClientOrder
	aux := &struct {
		OrderDate string `json:"orderDate"`
		*Alias
	}{
		Alias: (*Alias)(order),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return fmt.Errorf("invalid JSON structure: %v", err)
	}

	fmt.Println("Raw orderDate from JSON:", aux.OrderDate)

	parsedTime, err := time.Parse(time.RFC3339, aux.OrderDate)
	if err != nil {
		return fmt.Errorf("invalid date format for orderDate: %v", err)
	}

	order.OrderDate = parsedTime
	return nil
}
