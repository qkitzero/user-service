package user

import (
	"time"

	"github.com/qkitzero/user/internal/domain/user"
)

type UserUsecase interface {
	CreateUser(userID, displayName string, birthYear, birthMonth, birthDay int32) (user.User, error)
	GetUser(userID string) (user.User, error)
	UpdateUser(userID, displayName string) (user.User, error)
}

type userUsecase struct {
	repo user.UserRepository
}

func NewUserUsecase(repo user.UserRepository) UserUsecase {
	return &userUsecase{repo: repo}
}

func (s *userUsecase) CreateUser(id, displayName string, birthYear, birthMonth, birthDay int32) (user.User, error) {
	userID, err := user.NewUserID(id)
	if err != nil {
		return nil, err
	}

	userDisplayName, err := user.NewDisplayName(displayName)
	if err != nil {
		return nil, err
	}

	userBirthDate, err := user.NewBirthDate(birthYear, birthMonth, birthDay)
	if err != nil {
		return nil, err
	}

	user := user.NewUser(userID, userDisplayName, userBirthDate, time.Now(), time.Now())

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userUsecase) GetUser(id string) (user.User, error) {
	userID, err := user.NewUserID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.Read(userID)
}

func (s *userUsecase) UpdateUser(id, displayName string) (user.User, error) {
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
