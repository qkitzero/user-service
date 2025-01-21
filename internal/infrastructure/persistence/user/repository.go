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
	userTable := UserTable{
		ID:          user.ID(),
		DisplayName: user.DisplayName(),
		CreatedAt:   user.CreatedAt(),
		UpdatedAt:   user.UpdatedAt(),
	}

	if err := r.db.Create(&userTable).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Read(userID user.UserID) (user.User, error) {
	var userTable UserTable

	if err := r.db.First(&userTable, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return user.NewUser(
		userTable.ID,
		userTable.DisplayName,
		userTable.CreatedAt,
		userTable.UpdatedAt,
	), nil
}

func (r *userRepository) Update(user user.User) error {
	userTable := UserTable{
		ID:          user.ID(),
		DisplayName: user.DisplayName(),
		CreatedAt:   user.CreatedAt(),
		UpdatedAt:   user.UpdatedAt(),
	}

	if err := r.db.Save(&userTable).Error; err != nil {
		return err
	}

	return nil
}
