package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(dbHost, dbUser, dbPassword, dbName, dbPort, sslMode string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
