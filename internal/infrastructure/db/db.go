package db

import (
	"user/internal/infrastructure/persistence/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() (*gorm.DB, error) {
	dsn := "user:password@tcp(user-db:3306)/user_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&user.UserTable{})

	return db, nil
}
