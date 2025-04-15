package user

import (
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	mocks "github.com/qkitzero/user/mocks/domain/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		userID      string
		displayName string
		year        int32
		month       int32
		day         int32
		createErr   error
	}{
		{"success create user", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", 2000, 1, 1, nil},
		{"failure invalid user id", false, "0123456789", "test user", 2000, 1, 1, nil},
		{"failure empty display name", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", 2000, 1, 1, nil},
		{"failure birth date is in the future", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", 2500, 1, 1, nil},
		{"failure create error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", 2000, 1, 1, fmt.Errorf("create error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			userUsecase := NewUserUsecase(mockUserRepository)
			mockUserRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr).AnyTimes()
			_, err := userUsecase.CreateUser(tt.userID, tt.displayName, tt.year, tt.month, tt.day)
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
		name    string
		success bool
		userID  string
		readErr error
	}{
		{"success read user", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", nil},
		{"failure invalid user id", false, "0123456789", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			userUsecase := NewUserUsecase(mockUserRepository)
			mockUser := mocks.NewMockUser(ctrl)
			mockUserRepository.EXPECT().Read(gomock.Any()).Return(mockUser, nil).AnyTimes()
			_, err := userUsecase.GetUser(tt.userID)
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
		readErr     error
		updateErr   error
	}{
		{"success read user", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", nil, nil},
		{"failure invalid user id", false, "0123456789", "test user", nil, nil},
		{"failure read error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", fmt.Errorf("read error"), nil},
		{"failure empty display name", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", nil, nil},
		{"failure update error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", nil, fmt.Errorf("update error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUser := mocks.NewMockUser(ctrl)
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			userUsecase := NewUserUsecase(mockUserRepository)
			mockUserRepository.EXPECT().Read(gomock.Any()).Return(mockUser, tt.readErr).AnyTimes()
			mockUser.EXPECT().Update(gomock.Any()).AnyTimes()
			mockUserRepository.EXPECT().Update(gomock.Any()).Return(tt.updateErr).AnyTimes()
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
