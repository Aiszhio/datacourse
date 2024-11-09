package handlers

import (
	"fmt"
	"github.com/Aiszhio/datacourse.git/pkg/Redis"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetUserByPhone(rdb *redis.Client, db *gorm.DB, role string) (interface{}, error) {
	userInfo := make(map[string]interface{})
	var tableName string
	if role == "client" {
		tableName = "clients"
	} else {
		tableName = "employees"
	}

	phone, err := Redis.GetKey(rdb, "phone")
	if err != nil {
		return nil, fmt.Errorf("error in getting phone number from Redis: %w", err)
	}

	err = db.Table(tableName).Where("phone_number = ?", phone).Limit(1).Find(&userInfo).Error
	if err != nil {
		return nil, fmt.Errorf("error in getting user info from table %s: %w", tableName, err)
	}

	return userInfo, nil
}
