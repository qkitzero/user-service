package user

import (
	"time"

	"github.com/qkitzero/user/internal/domain/user"
)

type UserService interface {
	CreateUser(userID, displayName string) (user.User, error)
	GetUser(userID string) (user.User, error)
	UpdateUser(userID, displayName string) (user.User, error)
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(id, displayName string) (user.User, error) {
	userID, err := user.NewUserID(id)
	if err != nil {
		return nil, err
	}

	userDisplayName, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}

	user := user.NewUser(userID, userDisplayName, time.Now(), time.Now())
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(id string) (user.User, error) {
	userID, err := user.NewUserID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.Read(userID)
}

func (s *userService) UpdateUser(id, displayName string) (user.User, error) {
	userID, err := user.NewUserID(id)
	if err != nil {
		return nil, err
	}

	existingUser, err := s.repo.Read(userID)
	if err != nil {
		return nil, err
	}

	newDisplayName, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}

	existingUser.Update(newDisplayName)

	if err := s.repo.Update(existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}
