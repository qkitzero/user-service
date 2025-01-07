package user

import "user/internal/domain/user"

type UserService interface {
	CreateUser(displayName, email string) (user.User, error)
}

type userService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(displayName, email string) (user.User, error) {
	userID := user.NewUserID()

	userDisplayName, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}

	userEmail, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}

	user := user.NewUser(userID, userDisplayName, userEmail)
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
