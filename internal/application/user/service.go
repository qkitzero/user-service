package user

import "user/internal/domain/user"

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(email string) error {
	e, err := user.NewEmail(email)
	if err != nil {
		return err
	}
	user := user.NewUser(user.NewUserID(), e)
	return s.repo.Create(user)
}
