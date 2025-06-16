package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(dbUser, dbPassword, dbHost, dbPort, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
