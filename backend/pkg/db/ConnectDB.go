package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDB() (*gorm.DB, error) {
	dbLink := "postgres://postgres:1234@postgres:5432/parlourdb?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dbLink), &gorm.Config{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}
