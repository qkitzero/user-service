package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/qkitzero/auth/pb"
	"github.com/qkitzero/user/internal/application/auth"
	"google.golang.org/grpc/metadata"
)

type authUsecase struct {
	client pb.AuthServiceClient
}

func NewAuthUsecase(client pb.AuthServiceClient) auth.AuthUsecase {
	return &authUsecase{client: client}
}

func (s *authUsecase) VerifyToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata in context")
	}

	authHeader := md["authorization"]

	if len(authHeader) == 0 {
		return "", fmt.Errorf("authorization header is missing")
	}

	accessToken := strings.TrimPrefix(authHeader[0], "Bearer ")

	verifyTokenRequest := &pb.VerifyTokenRequest{AccessToken: accessToken}

	verifyTokenResponse, err := s.client.VerifyToken(ctx, verifyTokenRequest)
	if err != nil {
		return "", err
	}

	return verifyTokenResponse.GetUserId(), nil
}
