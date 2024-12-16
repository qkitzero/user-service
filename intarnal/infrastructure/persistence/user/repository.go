package user

import (
	"database/sql"
	"fmt"
	"user/intarnal/domain/user"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(user user.User) error {
	fmt.Println(user)
	return nil
}
