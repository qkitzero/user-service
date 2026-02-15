package user

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"

	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksappuser "github.com/qkitzero/user-service/mocks/application/user"
	mocksuser "github.com/qkitzero/user-service/mocks/domain/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		success       bool
		ctx           context.Context
		displayName   string
		year          int32
		month         int32
		day           int32
		createUserErr error
	}{
		{"success create user", true, context.Background(), "test user", 2000, 1, 1, nil},
		{"failure create user error", false, context.Background(), "test user", 2000, 1, 1, fmt.Errorf("create user error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			mockUserUsecase.EXPECT().CreateUser(tt.ctx, tt.displayName, tt.year, tt.month, tt.day).Return(mockUser, tt.createUserErr).AnyTimes()
			mockUserID := user.NewUserID()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()

			userHandler := NewUserHandler(mockUserUsecase)

			req := &userv1.CreateUserRequest{
				DisplayName: tt.displayName,
				BirthDate: &date.Date{
					Year:  tt.year,
					Month: tt.month,
					Day:   tt.day,
				},
			}

			_, err := userHandler.CreateUser(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		success    bool
		ctx        context.Context
		getUserErr error
	}{
		{"success get user", true, context.Background(), nil},
		{"failure get user error", false, context.Background(), fmt.Errorf("get user error")},
		{"failure user not found error", false, context.Background(), user.ErrUserNotFound},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			mockUserID := user.NewUserID()
			mockDisplayName, _ := user.NewDisplayName("test user")
			mockBirthDate, _ := user.NewBirthDate(2000, 1, 1)
			mockUserUsecase.EXPECT().GetUser(tt.ctx).Return(mockUser, tt.getUserErr).AnyTimes()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(mockDisplayName).AnyTimes()
			mockUser.EXPECT().BirthDate().Return(mockBirthDate).AnyTimes()

			userHandler := NewUserHandler(mockUserUsecase)

			req := &userv1.GetUserRequest{}

			_, err := userHandler.GetUser(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		success       bool
		ctx           context.Context
		displayName   string
		year          int32
		month         int32
		day           int32
		updateUserErr error
	}{
		{"success update user", true, context.Background(), "updated test user", 2000, 1, 1, nil},
		{"failure update user error", false, context.Background(), "updated test user", 2000, 1, 1, fmt.Errorf("update user error")},
		{"failure user not found error", false, context.Background(), "updated test user", 2000, 1, 1, user.ErrUserNotFound},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			mockUserID, _ := user.NewUserIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
			mockDisplayName, _ := user.NewDisplayName("test user")
			mockBirthDate, _ := user.NewBirthDate(2000, 1, 1)
			mockUserUsecase.EXPECT().UpdateUser(tt.ctx, tt.displayName, tt.year, tt.month, tt.day).Return(mockUser, tt.updateUserErr).AnyTimes()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(mockDisplayName).AnyTimes()
			mockUser.EXPECT().BirthDate().Return(mockBirthDate).AnyTimes()

			userHandler := NewUserHandler(mockUserUsecase)

			req := &userv1.UpdateUserRequest{
				DisplayName: tt.displayName,
				BirthDate: &date.Date{
					Year:  tt.year,
					Month: tt.month,
					Day:   tt.day,
				},
			}

			_, err := userHandler.UpdateUser(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
