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
		name         string
		success      bool
		userID       string
		displayName  string
		year         int32
		month        int32
		day          int32
		expectCreate bool
		createErr    error
	}{
		{"success create user", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", 2000, 1, 1, true, nil},
		{"failure invalid user id", false, "0123456789", "test user", 2000, 1, 1, false, nil},
		{"failure empty display name", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", 2000, 1, 1, false, nil},
		{"failure birth date is in the future", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", 2500, 1, 1, false, nil},
		{"failure create error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", 2000, 1, 1, true, fmt.Errorf("create error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			userService := NewUserService(mockUserRepository)
			if tt.expectCreate {
				mockUserRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr)
			}
			_, err := userService.CreateUser(tt.userID, tt.displayName, tt.year, tt.month, tt.day)
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
		name       string
		success    bool
		userID     string
		expectRead bool
		readErr    error
	}{
		{"success read user", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", true, nil},
		{"failure invalid user id", false, "0123456789", false, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			userService := NewUserService(mockUserRepository)
			mockUser := mocks.NewMockUser(ctrl)
			if tt.expectRead {
				mockUserRepository.EXPECT().Read(gomock.Any()).Return(mockUser, nil)
			}
			_, err := userService.GetUser(tt.userID)
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
		name         string
		success      bool
		userID       string
		displayName  string
		expectRead   bool
		readErr      error
		expectUpdate bool
		updateErr    error
	}{
		{"success read user", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", true, nil, true, nil},
		{"failure invalid user id", false, "0123456789", "test user", false, nil, false, nil},
		{"failure read error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", true, fmt.Errorf("read error"), false, nil},
		{"failure empty display name", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "", true, nil, false, nil},
		{"failure update error", false, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425", "test user", true, nil, true, fmt.Errorf("update error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUser := mocks.NewMockUser(ctrl)
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			userService := NewUserService(mockUserRepository)
			if tt.expectRead {
				mockUserRepository.EXPECT().Read(gomock.Any()).Return(mockUser, tt.readErr)
			}
			if tt.expectUpdate {
				mockUser.EXPECT().Update(gomock.Any())
			}
			if tt.expectUpdate {
				mockUserRepository.EXPECT().Update(gomock.Any()).Return(tt.updateErr)
			}
			_, err := userService.UpdateUser(tt.userID, tt.displayName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
