package membership

import (
	"context"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type MembershipRepository interface {
	Create(ctx context.Context, membership Membership) error
	FindByUserAndGroup(ctx context.Context, userID user.UserID, groupID group.GroupID) (Membership, error)
	ListByGroupID(ctx context.Context, groupID group.GroupID) ([]Membership, error)
	ListByUserID(ctx context.Context, userID user.UserID) ([]Membership, error)
	Update(ctx context.Context, membership Membership) error
	Delete(ctx context.Context, userID user.UserID, groupID group.GroupID) error
	CountOwners(ctx context.Context, groupID group.GroupID) (int, error)
}
