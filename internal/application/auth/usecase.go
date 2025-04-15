package auth

import "context"

type AuthUsecase interface {
	VerifyToken(ctx context.Context) (string, error)
}
