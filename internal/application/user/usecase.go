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
	CreateUser(ctx context.Context, displayName user.DisplayName, birthDate user.BirthDate) (user.User, error)
	GetUser(ctx context.Context) (user.User, error)
	UpdateUser(ctx context.Context, displayName user.DisplayName, birthDate user.BirthDate) (user.User, error)
}

type userUsecase struct {
	authService auth.AuthService
	userRepo    user.UserRepository
}

func NewUserUsecase(authService auth.AuthService, userRepo user.UserRepository) UserUsecase {
	return &userUsecase{authService: authService, userRepo: userRepo}
}

func (u *userUsecase) CreateUser(ctx context.Context, displayName user.DisplayName, birthDate user.BirthDate) (user.User, error) {
	identityID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	newIdentityID, err := identity.NewIdentityID(identityID)
	if err != nil {
		return nil, err
	}

	identities := []identity.Identity{identity.NewIdentity(newIdentityID)}

	now := time.Now()

	newUser := user.NewUser(user.NewUserID(), identities, displayName, birthDate, now, now)

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

func (u *userUsecase) UpdateUser(ctx context.Context, displayName user.DisplayName, birthDate user.BirthDate) (user.User, error) {
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

	foundUser.Update(displayName, birthDate)

	if err := u.userRepo.Update(foundUser); err != nil {
		return nil, err
	}

	return foundUser, nil
}
