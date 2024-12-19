package user

import "user/internal/domain/user"

type UserService struct {
	repo user.UserRepository
}

func NewUserService(repo user.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(displayName, email string) (user.User, error) {
	id := user.NewUserID()
	dn, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}
	e, err := user.NewEmail(email)
	if err != nil {
		return nil, err
	}
	user := user.NewUser(id, dn, e)
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}
