package user

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksappauth "github.com/qkitzero/user-service/mocks/application/auth"
	mocksuser "github.com/qkitzero/user-service/mocks/domain/user"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()
	displayName, _ := user.NewDisplayName("test user")
	birthDate, _ := user.NewBirthDate(2000, 1, 1)

	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		identityID     string
		verifyTokenErr error
		createErr      error
	}{
		{"success create user", true, context.Background(), "google-oauth2|000000000000000000000", nil, nil},
		{"failure verify token error", false, context.Background(), "", fmt.Errorf("verify token error"), nil},
		{"failure empty identity id", false, context.Background(), "", nil, nil},
		{"failure create error", false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepository := mocksuser.NewMockUserRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserRepository.EXPECT().Create(tt.ctx, gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewUserUsecase(mockAuthService, mockUserRepository)

			_, err := u.CreateUser(tt.ctx, displayName, birthDate)
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
			mockUser := mocksuser.NewMockUser(ctrl)
			mockUserRepository := mocksuser.NewMockUserRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserRepository.EXPECT().FindByIdentityID(tt.ctx, gomock.Any()).Return(mockUser, tt.findByIdentityIDErr).AnyTimes()

			u := NewUserUsecase(mockAuthService, mockUserRepository)

			_, err := u.GetUser(tt.ctx)
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
	displayName, _ := user.NewDisplayName("test user")
	birthDate, _ := user.NewBirthDate(2000, 1, 1)

	tests := []struct {
		name                string
		success             bool
		ctx                 context.Context
		identityID          string
		verifyTokenErr      error
		findByIdentityIDErr error
		updateErr           error
	}{
		{"success update user", true, context.Background(), "google-oauth2|000000000000000000000", nil, nil, nil},
		{"failure verify token error", false, context.Background(), "", fmt.Errorf("verify token error"), nil, nil},
		{"failure empty identity id", false, context.Background(), "", nil, nil, nil},
		{"failure find by identity id error", false, context.Background(), "google-oauth2|000000000000000000000", nil, errors.New("find by identity id error"), nil},
		{"failure identity not found", false, context.Background(), "google-oauth2|000000000000000000000", nil, identity.ErrIdentityNotFound, nil},
		{"failure update error", false, context.Background(), "google-oauth2|000000000000000000000", nil, nil, errors.New("update error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocksappauth.NewMockAuthService(ctrl)
			mockUser := mocksuser.NewMockUser(ctrl)
			mockUser.EXPECT().Update(gomock.Any(), gomock.Any()).AnyTimes()
			mockUserRepository := mocksuser.NewMockUserRepository(ctrl)
			mockAuthService.EXPECT().VerifyToken(tt.ctx).Return(tt.identityID, tt.verifyTokenErr).AnyTimes()
			mockUserRepository.EXPECT().FindByIdentityID(tt.ctx, gomock.Any()).Return(mockUser, tt.findByIdentityIDErr).AnyTimes()
			mockUserRepository.EXPECT().Update(tt.ctx, gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewUserUsecase(mockAuthService, mockUserRepository)

			_, err := u.UpdateUser(tt.ctx, displayName, birthDate)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
