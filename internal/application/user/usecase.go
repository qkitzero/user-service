package user

import (
	"errors"
	"time"

	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type UserUsecase interface {
	CreateUser(identityIDStr, displayNameStr string, y, m, d int32) (user.User, error)
	GetUser(identityIDStr string) (user.User, error)
	UpdateUser(userIDStr, displayNameStr string) (user.User, error)
}

type userUsecase struct {
	repo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (s *userUsecase) CreateUser(identityIDStr, displayNameStr string, y, m, d int32) (user.User, error) {
	identityID, err := identity.NewIdentityID(identityIDStr)
	if err != nil {
		return nil, err
	}

	identities := []identity.Identity{identity.NewIdentity(identityID)}

	displayName, err := user.NewDisplayName(displayNameStr)
	if err != nil {
		return nil, err
	}

	birthDate, err := user.NewBirthDate(y, m, d)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	u := user.NewUser(user.NewUserID(), identities, displayName, birthDate, now, now)

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userUsecase) GetUser(identityIDStr string) (user.User, error) {
	id, err := identity.NewIdentityID(identityIDStr)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.FindByIdentityID(id)
	if errors.Is(err, identity.ErrIdentityNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userUsecase) UpdateUser(userIDStr, displayNameStr string) (user.User, error) {
	userID, err := user.NewUserIDFromString(userIDStr)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	displayName, err := user.NewDisplayName(displayNameStr)
	if err != nil {
		return nil, err
	}

	u.Update(displayName)

	if err := s.repo.Update(u); err != nil {
		return nil, err
	}

	return u, nil
}
