package user

import (
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/user/internal/domain/identity"
	"github.com/qkitzero/user/internal/domain/user"
	infraidentity "github.com/qkitzero/user/internal/infrastructure/identity"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(u user.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		userModel := UserModel{
			ID:          u.ID(),
			DisplayName: u.DisplayName(),
			BirthDate:   u.BirthDate(),
			CreatedAt:   u.CreatedAt(),
			UpdatedAt:   u.UpdatedAt(),
		}

		if err := tx.Create(&userModel).Error; err != nil {
			return err
		}

		var identityModels []infraidentity.IdentityModel
		for _, i := range u.Identities() {
			identityModels = append(identityModels, infraidentity.IdentityModel{
				ID:     i.ID(),
				UserID: u.ID(),
			})
		}

		if err := tx.Create(&identityModels).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *userRepository) FindByIdentityID(identityID identity.IdentityID) (user.User, error) {
	var identityModel infraidentity.IdentityModel
	err := r.db.Where("id = ?", identityID).First(&identityModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, identity.ErrIdentityNotFound
	}
	if err != nil {
		return nil, err
	}

	var userModel UserModel
	err = r.db.Where("id = ?", identityModel.UserID).First(&userModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	var identityModels []infraidentity.IdentityModel
	if err := r.db.Where("user_id = ?", userModel.ID).Find(&identityModels).Error; err != nil {
		return nil, err
	}

	var identities []identity.Identity
	for _, im := range identityModels {
		identities = append(identities, identity.NewIdentity(im.ID))
	}

	return user.NewUser(
		userModel.ID,
		identities,
		userModel.DisplayName,
		userModel.BirthDate,
		userModel.CreatedAt,
		userModel.UpdatedAt,
	), nil
}

func (r *userRepository) FindByID(userID user.UserID) (user.User, error) {
	var userModel UserModel
	err := r.db.Where("id = ?", userID).First(&userModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	var identityModels []infraidentity.IdentityModel
	if err := r.db.Where("user_id = ?", userModel.ID).Find(&identityModels).Error; err != nil {
		return nil, err
	}

	var identities []identity.Identity
	for _, im := range identityModels {
		identities = append(identities, identity.NewIdentity(im.ID))
	}

	return user.NewUser(
		userModel.ID,
		identities,
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
