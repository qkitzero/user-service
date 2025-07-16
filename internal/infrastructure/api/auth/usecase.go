package auth

import (
	"context"
	"errors"

	authv1 "github.com/qkitzero/auth/gen/go/auth/v1"
	"github.com/qkitzero/user-service/internal/application/auth"
	"google.golang.org/grpc/metadata"
)

type authUsecase struct {
	client authv1.AuthServiceClient
}

func NewAuthUsecase(client authv1.AuthServiceClient) auth.AuthUsecase {
	return &authUsecase{client: client}
}

func (s *authUsecase) VerifyToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is missing")
	}

	ctx = metadata.NewOutgoingContext(ctx, md)

	verifyTokenRequest := &authv1.VerifyTokenRequest{}

	verifyTokenResponse, err := s.client.VerifyToken(ctx, verifyTokenRequest)
	if err != nil {
		return "", err
	}

	return verifyTokenResponse.GetUserId(), nil
}
