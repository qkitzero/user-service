package user

import (
	"context"
	"errors"
	"time"

	"github.com/qkitzero/user-service/internal/application/auth"
	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, displayName string, year, month, day int32) (user.User, error)
	GetUser(ctx context.Context) (user.User, error)
	UpdateUser(ctx context.Context, displayName string, year, month, day int32) (user.User, error)
}

type userUsecase struct {
	authService auth.AuthService
	userRepo    user.UserRepository
}

func NewUserUsecase(authService auth.AuthService, userRepo user.UserRepository) UserUsecase {
	return &userUsecase{authService: authService, userRepo: userRepo}
}

func (u *userUsecase) CreateUser(ctx context.Context, displayName string, year, month, day int32) (user.User, error) {
	identityID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	newIdentityID, err := identity.NewIdentityID(identityID)
	if err != nil {
		return nil, err
	}

	identities := []identity.Identity{identity.NewIdentity(newIdentityID)}

	newDisplayName, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}

	newBirthDate, err := user.NewBirthDate(year, month, day)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	newUser := user.NewUser(user.NewUserID(), identities, newDisplayName, newBirthDate, now, now)

	if err := u.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (u *userUsecase) GetUser(ctx context.Context) (user.User, error) {
	identityID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	id, err := identity.NewIdentityID(identityID)
	if err != nil {
		return nil, err
	}

	foundUser, err := u.userRepo.FindByIdentityID(id)
	if errors.Is(err, identity.ErrIdentityNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (u *userUsecase) UpdateUser(ctx context.Context, displayName string, year, month, day int32) (user.User, error) {
	identityID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	id, err := identity.NewIdentityID(identityID)
	if err != nil {
		return nil, err
	}

	foundUser, err := u.userRepo.FindByIdentityID(id)
	if errors.Is(err, identity.ErrIdentityNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	newDisplayName, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}

	newBirthDate, err := user.NewBirthDate(year, month, day)
	if err != nil {
		return nil, err
	}

	foundUser.Update(newDisplayName, newBirthDate)

	if err := u.userRepo.Update(foundUser); err != nil {
		return nil, err
	}

	return foundUser, nil
}
