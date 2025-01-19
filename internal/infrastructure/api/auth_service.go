package api

import (
	"context"
	"fmt"
	"strings"
	"user/internal/application/auth"

	"github.com/qkitzero/auth/pb"
	"google.golang.org/grpc/metadata"
)

type authService struct {
	client pb.AuthServiceClient
}

func NewAuthService(client pb.AuthServiceClient) auth.AuthService {
	return &authService{client: client}
}

func (s *authService) VerifyToken(ctx context.Context) (string, error) {
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
