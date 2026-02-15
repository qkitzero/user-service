package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/user-service/internal/domain/identity"
	mocksappauth "github.com/qkitzero/user-service/mocks/application/auth"
	mocks "github.com/qkitzero/user-service/mocks/domain/user"
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
		createErr      error
	}{
		{"success create user", true, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, nil},
		{"failure verify token error", false, context.Background(), "", "test user", 2000, 1, 1, fmt.Errorf("verify token error"), nil},
		{"failure empty identity id", false, context.Background(), "", "test user", 2000, 1, 1, nil, nil},
		{"failure empty display name", false, context.Background(), "google-oauth2|000000000000000000000", "", 2000, 1, 1, nil, nil},
		{"failure birth date is in the future", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2500, 1, 1, nil, nil},
		{"failure create error", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserRepository.EXPECT().Create(gomock.Any()).Return(tt.createErr).AnyTimes()

			userUsecase := NewUserUsecase(mockAuthService, mockUserRepository)

			_, err := userUsecase.CreateUser(tt.ctx, tt.displayName, tt.year, tt.month, tt.day)
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
		name                string
		success             bool
		ctx                 context.Context
		identityID          string
		verifyTokenErr      error
		findByIdentityIDErr error
	}{
		{"success get user", true, context.Background(), "google-oauth2|000000000000000000000", nil, nil},
		{"failure verify token error", false, context.Background(), "", fmt.Errorf("verify token error"), nil},
		{"failure empty identity id", false, context.Background(), "", nil, nil},
		{"failure find by identity id error", false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by identity id error")},
		{"failure identity not found", false, context.Background(), "google-oauth2|000000000000000000000", nil, identity.ErrIdentityNotFound},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUser := mocks.NewMockUser(ctrl)
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserRepository.EXPECT().FindByIdentityID(gomock.Any()).Return(mockUser, tt.findByIdentityIDErr).AnyTimes()

			userUsecase := NewUserUsecase(mockAuthService, mockUserRepository)

			_, err := userUsecase.GetUser(tt.ctx)
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
		name                string
		success             bool
		ctx                 context.Context
		identityID          string
		displayName         string
		year                int32
		month               int32
		day                 int32
		verifyTokenErr      error
		findByIdentityIDErr error
		updateErr           error
	}{
		{"success update user", true, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, nil, nil},
		{"failure verify token error", false, context.Background(), "", "test user", 2000, 1, 1, fmt.Errorf("verify token error"), nil, nil},
		{"failure empty identity id", false, context.Background(), "", "test user", 2000, 1, 1, nil, nil, nil},
		{"failure find by identity id error", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, errors.New("find by identity id error"), nil},
		{"failure identity not found", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, identity.ErrIdentityNotFound, nil},
		{"failure empty display name", false, context.Background(), "google-oauth2|000000000000000000000", "", 2000, 1, 1, nil, nil, nil},
		{"failure invalid birth date", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 2, 30, nil, nil, nil},
		{"failure update error", false, context.Background(), "google-oauth2|000000000000000000000", "test user", 2000, 1, 1, nil, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUser := mocks.NewMockUser(ctrl)
			mockUser.EXPECT().Update(gomock.Any(), gomock.Any()).AnyTimes()
			mockUserRepository := mocks.NewMockUserRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserRepository.EXPECT().FindByIdentityID(gomock.Any()).Return(mockUser, tt.findByIdentityIDErr).AnyTimes()
			mockUserRepository.EXPECT().Update(gomock.Any()).Return(tt.updateErr).AnyTimes()

			userUsecase := NewUserUsecase(mockAuthService, mockUserRepository)

			_, err := userUsecase.UpdateUser(tt.ctx, tt.displayName, tt.year, tt.month, tt.day)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
