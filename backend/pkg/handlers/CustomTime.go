package handlers

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct CustomTime) ToTime() time.Time {
	return ct.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	str := string(b)

	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}

	str = strings.Replace(str, "T", " ", 1)

	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		t, err = time.Parse("2006-01-02", str)
		if err != nil {
			return fmt.Errorf("unable to parse time: %v", err)
		}
	}

	*ct = CustomTime{t}
	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(time.RFC3339))), nil
}

func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time.Format(time.RFC3339), nil
}

type CustomFloat64 float64

func (f *CustomFloat64) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return fmt.Errorf("invalid float64 format: %v", err)
	}
	*f = CustomFloat64(val)
	return nil
}
