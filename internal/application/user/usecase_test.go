package user

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/user-service/internal/domain/identity"
	mocks "github.com/qkitzero/user-service/mocks/domain/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		identityID  string
		displayName string
		year        int32
		month       int32
		day         int32
		createErr   error
	}{
		{"success create user", true, "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil},
		{"failure empty identity id", false, "", "test user", 2000, 1, 1, nil},
		{"failure empty display name", false, "google-oauth2|000000000000000000000", "", 2000, 1, 1, nil},
		{"failure birth date is in the future", false, "google-oauth2|000000000000000000000", "test user", 2500, 1, 1, nil},
		{"failure create error", false, "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			mockUserRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr).AnyTimes()

			userUsecase := NewUserUsecase(mockUserRepository)

			_, err := userUsecase.CreateUser(tt.identityID, tt.displayName, tt.year, tt.month, tt.day)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                string
		success             bool
		identityID          string
		findByIdentityIDErr error
	}{
		{"success get user", true, "google-oauth2|000000000000000000000", nil},
		{"failure empty identity id", false, "", nil},
		{"failure find by identity id error", false, "google-oauth2|000000000000000000000", errors.New("find by identity id error")},
		{"failure identity not found", false, "google-oauth2|000000000000000000000", identity.ErrIdentityNotFound},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := mocks.NewMockUser(ctrl)
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			mockUserRepository.EXPECT().FindByIdentityID(gomock.Any()).Return(mockUser, tt.findByIdentityIDErr).AnyTimes()

			userUsecase := NewUserUsecase(mockUserRepository)

			_, err := userUsecase.GetUser(tt.identityID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		userID      string
		displayName string
		findByIDErr error
		updateErr   error
	}{
		{"success update user", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", nil, nil},
		{"failure invalid user id", false, "0123456789", "test user", nil, nil},
		{"failure find by id error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", errors.New("find by id error"), nil},
		{"failure empty display name", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", nil, nil},
		{"failure update error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := mocks.NewMockUser(ctrl)
			mockUser.EXPECT().Update(gomock.Any()).AnyTimes()
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			mockUserRepository.EXPECT().FindByID(gomock.Any()).Return(mockUser, tt.findByIDErr).AnyTimes()
			mockUserRepository.EXPECT().Update(gomock.Any()).Return(tt.updateErr).AnyTimes()

			userUsecase := NewUserUsecase(mockUserRepository)

			_, err := userUsecase.UpdateUser(tt.userID, tt.displayName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
