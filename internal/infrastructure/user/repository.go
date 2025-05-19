package user

import (
	"github.com/qkitzero/user/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user user.User) error {
	userModel := UserModel{
		ID:          user.ID(),
		DisplayName: user.DisplayName(),
		BirthDate:   user.BirthDate(),
		CreatedAt:   user.CreatedAt(),
		UpdatedAt:   user.UpdatedAt(),
	}

	if err := r.db.Create(&userModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Read(userID user.UserID) (user.User, error) {
	var userModel UserModel

	if err := r.db.First(&userModel, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return user.NewUser(
		userModel.ID,
		userModel.DisplayName,
		userModel.BirthDate,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	), nil
}

func (r *userRepository) Update(user user.User) error {
	userModel := UserModel{
		ID:          user.ID(),
		DisplayName: user.DisplayName(),
		BirthDate:   user.BirthDate(),
		CreatedAt:   user.CreatedAt(),
		UpdatedAt:   user.UpdatedAt(),
	}

	if err := r.db.Save(&userModel).Error; err != nil {
		return err
	}

	return nil
}
