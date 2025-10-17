package user

import (
	"errors"
	"time"

	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type UserUsecase interface {
	CreateUser(identityID, displayName string, year, month, day int32) (user.User, error)
	GetUser(identityID string) (user.User, error)
	UpdateUser(identityID, displayName string, year, month, day int32) (user.User, error)
}

type userUsecase struct {
	repo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (s *userUsecase) CreateUser(identityID, displayName string, year, month, day int32) (user.User, error) {
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

	if err := s.repo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *userUsecase) GetUser(identityID string) (user.User, error) {
	id, err := identity.NewIdentityID(identityID)
	if err != nil {
		return nil, err
	}

	foundUser, err := s.repo.FindByIdentityID(id)
	if errors.Is(err, identity.ErrIdentityNotFound) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (s *userUsecase) UpdateUser(identityID, displayName string, year, month, day int32) (user.User, error) {
	id, err := identity.NewIdentityID(identityID)
	if err != nil {
		return nil, err
	}

	foundUser, err := s.repo.FindByIdentityID(id)
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

	if err := s.repo.Update(foundUser); err != nil {
		return nil, err
	}

	return foundUser, nil
}
