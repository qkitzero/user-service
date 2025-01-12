package user

import (
	"time"
	"user/internal/domain/user"
)

type UserService interface {
	CreateUser(displayName string) (user.User, error)
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(displayName string) (user.User, error) {
	userID := user.NewUserID()

	userDisplayName, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}

	user := user.NewUser(userID, userDisplayName, time.Now())
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
