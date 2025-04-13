package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/qkitzero/user/pb"
	"go.uber.org/mock/gomock"

	"github.com/qkitzero/user/internal/domain/user"
	mocksAuthService "github.com/qkitzero/user/mocks/application/auth"
	mocksUserService "github.com/qkitzero/user/mocks/application/user"
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
		{"success create user", true, context.Background(), "userID", "test user", 2000, 1, 1, nil, nil},
		{"failure verify token error", false, context.Background(), "userID", "test user", 2000, 1, 1, fmt.Errorf("verify token error"), nil},
		{"failure create user error", false, context.Background(), "userID", "test user", 2000, 1, 1, nil, fmt.Errorf("create user error")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAuthService := mocksAuthService.NewMockAuthService(ctrl)
			mockUserService := mocksUserService.NewMockUserService(ctrl)
			handler := NewUserHandler(mockAuthService, mockUserService)
			mockUser := mocksUser.NewMockUser(ctrl)
			mockUserID, _ := user.NewUserID(tt.userID)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.userID, tt.verifyTokenErr).AnyTimes()
			mockUserService.EXPECT().CreateUser(tt.userID, tt.displayName, tt.year, tt.month, tt.day).Return(mockUser, tt.createUserErr).AnyTimes()
			mockUser.EXPECT().ID().Return(mockUserID).AnyTimes()
			req := &pb.CreateUserRequest{
				DisplayName: tt.displayName,
				BirthDate: &pb.Date{
					Year:  tt.year,
					Month: tt.month,
					Day:   tt.day,
				},
			}
			_, err := handler.CreateUser(tt.ctx, req)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
