package user

import (
	"fmt"
	"user/intarnal/domain/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(user user.User) error {
	fmt.Println(user)
	return nil
}
