package gorm

import (
	"fmt"

	"gorm.io/driver/postgres"
	gormpkg "gorm.io/gorm"
)

func InitDB(host, dbName, user, password string, port int) *gormpkg.DB {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s port=%d sslmode=disable",
		host, dbName, user, password, port)

	db, err := gormpkg.Open(postgres.Open(dsn), &gormpkg.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных")
	}

	return db
}
