package auth

import "context"

type AuthService interface {
	VerifyToken(ctx context.Context) (string, error)
}
