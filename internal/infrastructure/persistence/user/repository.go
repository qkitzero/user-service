package user

import (
	"user/internal/domain/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user user.User) error {
	userTable := UserTable{
		ID:    user.ID(),
		Email: user.Email(),
	}
	r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&userTable).Error; err != nil {
			return err
		}
		return nil
	})
	return nil
}
