package user

import "user/intarnal/domain/user"

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(name, email string) error {
	user := user.NewUser("id", name, email)
	return s.repo.Create(user)
}
