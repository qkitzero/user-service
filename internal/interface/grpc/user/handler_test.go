package user

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksappuser "github.com/qkitzero/user-service/mocks/application/user"
	mocksuser "github.com/qkitzero/user-service/mocks/domain/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		ctx           context.Context
		displayName   string
		year          int32
		month         int32
		day           int32
		callUsecase   bool
		createUserErr error
		wantCode      codes.Code
	}{
		{"success create user", context.Background(), "test user", 2000, 1, 1, true, nil, codes.OK},
		{"failure invalid display name", context.Background(), "", 2000, 1, 1, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", context.Background(), "test user", 2000, 2, 30, false, nil, codes.InvalidArgument},
		{"failure usecase error", context.Background(), "test user", 2000, 1, 1, true, fmt.Errorf("create user error"), codes.Internal},
		{"failure status preserved", context.Background(), "test user", 2000, 1, 1, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().CreateUser(tt.ctx, gomock.Any(), gomock.Any()).Return(mockUser, tt.createUserErr).Times(1)
				mockUser.EXPECT().ID().Return(user.NewUserID()).AnyTimes()
			}

			handler := NewUserHandler(mockUsecase)

			req := &userv1.CreateUserRequest{
				DisplayName: tt.displayName,
				BirthDate: &date.Date{
					Year:  tt.year,
					Month: tt.month,
					Day:   tt.day,
				},
			}

			_, err := handler.CreateUser(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	t.Parallel()

	mockUserSample := func(ctrl *gomock.Controller) *mocksuser.MockUser {
		m := mocksuser.NewMockUser(ctrl)
		displayName, _ := user.NewDisplayName("test user")
		birthDate, _ := user.NewBirthDate(2000, 1, 1)
		m.EXPECT().ID().Return(user.NewUserID()).AnyTimes()
		m.EXPECT().DisplayName().Return(displayName).AnyTimes()
		m.EXPECT().BirthDate().Return(birthDate).AnyTimes()
		return m
	}

	tests := []struct {
		name       string
		ctx        context.Context
		getUserErr error
		wantCode   codes.Code
	}{
		{"success get user", context.Background(), nil, codes.OK},
		{"failure get user error", context.Background(), fmt.Errorf("get user error"), codes.Internal},
		{"failure user not found", context.Background(), user.ErrUserNotFound, codes.NotFound},
		{"failure status preserved", context.Background(), status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			mockUsecase.EXPECT().GetUser(tt.ctx).Return(mockUserSample(ctrl), tt.getUserErr).Times(1)

			handler := NewUserHandler(mockUsecase)

			_, err := handler.GetUser(tt.ctx, &userv1.GetUserRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()

	mockUserSample := func(ctrl *gomock.Controller) *mocksuser.MockUser {
		m := mocksuser.NewMockUser(ctrl)
		displayName, _ := user.NewDisplayName("test user")
		birthDate, _ := user.NewBirthDate(2000, 1, 1)
		m.EXPECT().ID().Return(user.NewUserID()).AnyTimes()
		m.EXPECT().DisplayName().Return(displayName).AnyTimes()
		m.EXPECT().BirthDate().Return(birthDate).AnyTimes()
		return m
	}

	tests := []struct {
		name          string
		ctx           context.Context
		displayName   string
		year          int32
		month         int32
		day           int32
		callUsecase   bool
		updateUserErr error
		wantCode      codes.Code
	}{
		{"success update user", context.Background(), "updated test user", 2000, 1, 1, true, nil, codes.OK},
		{"failure invalid display name", context.Background(), "", 2000, 1, 1, false, nil, codes.InvalidArgument},
		{"failure invalid birth date", context.Background(), "updated test user", 2000, 2, 30, false, nil, codes.InvalidArgument},
		{"failure usecase error", context.Background(), "updated test user", 2000, 1, 1, true, fmt.Errorf("update user error"), codes.Internal},
		{"failure user not found", context.Background(), "updated test user", 2000, 1, 1, true, user.ErrUserNotFound, codes.NotFound},
		{"failure status preserved", context.Background(), "updated test user", 2000, 1, 1, true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappuser.NewMockUserUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().UpdateUser(tt.ctx, gomock.Any(), gomock.Any()).Return(mockUserSample(ctrl), tt.updateUserErr).Times(1)
			}

			handler := NewUserHandler(mockUsecase)

			req := &userv1.UpdateUserRequest{
				DisplayName: tt.displayName,
				BirthDate: &date.Date{
					Year:  tt.year,
					Month: tt.month,
					Day:   tt.day,
				},
			}

			_, err := handler.UpdateUser(tt.ctx, req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
