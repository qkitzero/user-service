package group

import (
	"context"

	"github.com/qkitzero/user-service/internal/domain/user"
)

type GroupRepository interface {
	Create(ctx context.Context, group Group, ownerUserID user.UserID) error
	FindByID(ctx context.Context, groupID GroupID) (Group, error)
	FindByIDs(ctx context.Context, groupIDs []GroupID) ([]Group, error)
	Update(ctx context.Context, group Group) error
	Delete(ctx context.Context, groupID GroupID) error
	AddChild(ctx context.Context, parentID, childID GroupID) error
	RemoveChild(ctx context.Context, parentID, childID GroupID) error
	FindChildren(ctx context.Context, parentID GroupID) ([]Group, error)
	FindParents(ctx context.Context, childID GroupID) ([]Group, error)
}
