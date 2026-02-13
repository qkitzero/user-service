package auth

import (
	"context"
	"errors"
	"testing"

	authv1 "github.com/qkitzero/auth-service/gen/go/auth/v1"
	mocks "github.com/qkitzero/user-service/mocks/external/auth/v1"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

func TestVerifyToken(t *testing.T) {
	t.Parallel()
	accessToken := "accessToken"
	tests := []struct {
		name           string
		success        bool
		ctx            context.Context
		verifyTokenErr error
	}{
		{
			name:           "success verify token",
			success:        true,
			ctx:            metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+accessToken)),
			verifyTokenErr: nil,
		},
		{
			name:           "failure missing metadata",
			success:        false,
			ctx:            context.Background(),
			verifyTokenErr: nil,
		},
		{
			name:           "failure verify token error",
			success:        false,
			ctx:            metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+accessToken)),
			verifyTokenErr: errors.New("verify token error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockAuthServiceClient(ctrl)
			mockVerifyTokenResponse := &authv1.VerifyTokenResponse{
				UserId: "google-oauth2|000000000000000000000",
			}
			mockClient.EXPECT().VerifyToken(gomock.Any(), gomock.Any()).Return(mockVerifyTokenResponse, tt.verifyTokenErr).AnyTimes()

			authService := NewAuthService(mockClient)

			_, err := authService.VerifyToken(tt.ctx)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
