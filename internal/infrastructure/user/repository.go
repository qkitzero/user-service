package user

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/user"
	infraidentity "github.com/qkitzero/user-service/internal/infrastructure/identity"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &userRepository{db: db}
}

func toUser(m UserModel, identities []identity.Identity) user.User {
	return user.NewUser(m.ID, identities, m.DisplayName, m.BirthDate, m.CreatedAt, m.UpdatedAt)
}

func (r *userRepository) Create(ctx context.Context, u user.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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

func (r *userRepository) FindByIdentityID(ctx context.Context, id identity.IdentityID) (user.User, error) {
	var identityModel infraidentity.IdentityModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&identityModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, identity.ErrIdentityNotFound
	}
	if err != nil {
		return nil, err
	}

	var userModel UserModel
	err = r.db.WithContext(ctx).Where("id = ?", identityModel.UserID).First(&userModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	var identityModels []infraidentity.IdentityModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userModel.ID).Find(&identityModels).Error; err != nil {
		return nil, err
	}

	var identities []identity.Identity
	for _, im := range identityModels {
		identities = append(identities, identity.NewIdentity(im.ID))
	}

	return toUser(userModel, identities), nil
}

func (r *userRepository) FindByID(ctx context.Context, id user.UserID) (user.User, error) {
	var userModel UserModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	var identityModels []infraidentity.IdentityModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userModel.ID).Find(&identityModels).Error; err != nil {
		return nil, err
	}

	var identities []identity.Identity
	for _, im := range identityModels {
		identities = append(identities, identity.NewIdentity(im.ID))
	}

	return toUser(userModel, identities), nil
}

func (r *userRepository) FindByIDs(ctx context.Context, ids []user.UserID) ([]user.User, error) {
	if len(ids) == 0 {
		return []user.User{}, nil
	}

	var userModels []UserModel
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&userModels).Error; err != nil {
		return nil, err
	}

	if len(userModels) == 0 {
		return []user.User{}, nil
	}

	userIDs := make([]user.UserID, 0, len(userModels))
	for _, m := range userModels {
		userIDs = append(userIDs, m.ID)
	}

	var identityModels []infraidentity.IdentityModel
	if err := r.db.WithContext(ctx).Where("user_id IN ?", userIDs).Find(&identityModels).Error; err != nil {
		return nil, err
	}

	identitiesByUserID := make(map[user.UserID][]identity.Identity, len(userModels))
	for _, im := range identityModels {
		identitiesByUserID[im.UserID] = append(identitiesByUserID[im.UserID], identity.NewIdentity(im.ID))
	}

	users := make([]user.User, 0, len(userModels))
	for _, m := range userModels {
		users = append(users, toUser(m, identitiesByUserID[m.ID]))
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, u user.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userModel := UserModel{
			ID:          u.ID(),
			DisplayName: u.DisplayName(),
			BirthDate:   u.BirthDate(),
			CreatedAt:   u.CreatedAt(),
			UpdatedAt:   u.UpdatedAt(),
		}

		if err := tx.Save(&userModel).Error; err != nil {
			return err
		}

		return nil
	})
}
