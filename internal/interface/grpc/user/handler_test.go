package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/qkitzero/user/pb"
	"go.uber.org/mock/gomock"

	"github.com/qkitzero/user/internal/domain/user"
	mocksAuthUsecase "github.com/qkitzero/user/mocks/application/auth"
	mocksUserUsecase "github.com/qkitzero/user/mocks/application/user"
	mocksUser "github.com/qkitzero/user/mocks/domain/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		userID         string
		displayName    string
		year           int32
		month          int32
		day            int32
		verifyTokenErr error
		createUserErr  error
	}{
		{"success create user", true, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", "test user", 2000, 1, 1, nil, nil},
		{"failure verify token error", false, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", "test user", 2000, 1, 1, fmt.Errorf("verify token error"), nil},
		{"failure create user error", false, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", "test user", 2000, 1, 1, nil, fmt.Errorf("create user error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mocksAuthUsecase.NewMockAuthUsecase(ctrl)
			mockUserUsecase := mocksUserUsecase.NewMockUserUsecase(ctrl)
			mockUser := mocksUser.NewMockUser(ctrl)
			mockAuthUsecase.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserUsecase.EXPECT().CreateUser(tt.userID, tt.displayName, tt.year, tt.month, tt.day).Return(mockUser, tt.createUserErr).AnyTimes()
			mockUserID, _ := user.NewUserID(tt.userID)
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()

			userHandler := NewUserHandler(mockAuthUsecase, mockUserUsecase)

			req := &pb.CreateUserRequest{
				DisplayName: tt.displayName,
				BirthDate: &pb.Date{
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
		userID         string
		verifyTokenErr error
		getUserErr     error
	}{
		{"success get user", true, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", nil, nil},
		{"failure verify token error", false, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", fmt.Errorf("verify token error"), nil},
		{"failure get user error", false, context.Background(), "0800819d-746e-4cb9-b561-6841b98cb19c", nil, fmt.Errorf("get user error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mocksAuthUsecase.NewMockAuthUsecase(ctrl)
			mockUserUsecase := mocksUserUsecase.NewMockUserUsecase(ctrl)
			mockUser := mocksUser.NewMockUser(ctrl)
			mockUserID, _ := user.NewUserID(tt.userID)
			mockDisplayName, _ := user.NewDisplayName("test user")
			mockBirthDate, _ := user.NewBirthDate(2000, 1, 1)
			mockAuthUsecase.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserUsecase.EXPECT().GetUser(tt.userID).Return(mockUser, tt.getUserErr).AnyTimes()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(mockDisplayName).AnyTimes()
			mockUser.EXPECT().BirthDate().Return(mockBirthDate).AnyTimes()

			userHandler := NewUserHandler(mockAuthUsecase, mockUserUsecase)

			req := &pb.GetUserRequest{}

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
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthUsecase := mocksAuthUsecase.NewMockAuthUsecase(ctrl)
			mockUserUsecase := mocksUserUsecase.NewMockUserUsecase(ctrl)
			mockUser := mocksUser.NewMockUser(ctrl)
			mockUserID, _ := user.NewUserID(tt.userID)
			mockDisplayName, _ := user.NewDisplayName(tt.displayName)
			mockAuthUsecase.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserUsecase.EXPECT().UpdateUser(tt.userID, tt.displayName).Return(mockUser, tt.updateUserErr).AnyTimes()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()
			mockUser.EXPECT().DisplayName().Return(mockDisplayName).AnyTimes()

			userHandler := NewUserHandler(mockAuthUsecase, mockUserUsecase)

			req := &pb.UpdateUserRequest{
				DisplayName: tt.displayName,
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
