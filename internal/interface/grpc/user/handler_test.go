package user

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"

	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksappauth "github.com/qkitzero/user-service/mocks/application/auth"
	mocksappuser "github.com/qkitzero/user-service/mocks/application/user"
	mocksuser "github.com/qkitzero/user-service/mocks/domain/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		identityID     string
		displayName    string
		year           int32
		month          int32
		day            int32
		verifyTokenErr error
		createUserErr  error
	}{
		{"success create user", true, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, nil},
		{"failure verify token error", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, fmt.Errorf("verify token error"), nil},
		{"failure create user error", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, fmt.Errorf("create user error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mocksappauth.NewMockAuthUsecase(ctrl)
			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			mockAuthUsecase.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserUsecase.EXPECT().CreateUser(tt.identityID, tt.displayName, tt.year, tt.month, tt.day).Return(mockUser, tt.createUserErr).AnyTimes()
			mockUserID := user.NewUserID()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()

			userHandler := NewUserHandler(mockAuthUsecase, mockUserUsecase)

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
				t.Errorf("expected error but got nil")
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		identityID     string
		verifyTokenErr error
		getUserErr     error
	}{
		{"success get user", true, context.Background(), "google-oauth2|000000000000000000000", nil, nil},
		{"failure verify token error", false, context.Background(), "google-oauth2|000000000000000000000", fmt.Errorf("verify token error"), nil},
		{"failure get user error", false, context.Background(), "google-oauth2|000000000000000000000", nil, fmt.Errorf("get user error")},
		{"failure user not found error", false, context.Background(), "google-oauth2|000000000000000000000", nil, user.ErrUserNotFound},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mocksappauth.NewMockAuthUsecase(ctrl)
			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			mockUserID := user.NewUserID()
			mockDisplayName, _ := user.NewDisplayName("test user")
			mockBirthDate, _ := user.NewBirthDate(2000, 1, 1)
			mockAuthUsecase.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserUsecase.EXPECT().GetUser(tt.identityID).Return(mockUser, tt.getUserErr).AnyTimes()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(mockDisplayName).AnyTimes()
			mockUser.EXPECT().BirthDate().Return(mockBirthDate).AnyTimes()

			userHandler := NewUserHandler(mockAuthUsecase, mockUserUsecase)

			req := &userv1.GetUserRequest{}

			_, err := userHandler.GetUser(tt.ctx, req)
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
		name           string
		success        bool
		ctx            context.Context
		userID         string
		displayName    string
		verifyTokenErr error
		updateUserErr  error
	}{
		{"success update user", true, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", "test user", nil, nil},
		{"failure verify token error", false, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", "test user", fmt.Errorf("verify token error"), nil},
		{"failure update user error", false, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", "test user", nil, fmt.Errorf("update user error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mocksappauth.NewMockAuthUsecase(ctrl)
			mockUserUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			mockUserID, _ := user.NewUserIDFromString(tt.userID)
			mockDisplayName, _ := user.NewDisplayName(tt.displayName)
			mockAuthUsecase.EXPECT().VerifyToken(tt.ctx).Return("", tt.verifyTokenErr).AnyTimes()
			mockUserUsecase.EXPECT().UpdateUser(tt.userID, tt.displayName).Return(mockUser, tt.updateUserErr).AnyTimes()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(mockDisplayName).AnyTimes()

			userHandler := NewUserHandler(mockAuthUsecase, mockUserUsecase)

			req := &userv1.UpdateUserRequest{
				DisplayName: tt.displayName,
				UserId:      tt.userID,
			}

			_, err := userHandler.UpdateUser(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
