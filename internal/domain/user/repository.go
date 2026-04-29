package user

import (
	"context"

	"github.com/qkitzero/user-service/internal/domain/identity"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindByIdentityID(ctx context.Context, identityID identity.IdentityID) (User, error)
	FindByID(ctx context.Context, userID UserID) (User, error)
	Update(ctx context.Context, user User) error
}
