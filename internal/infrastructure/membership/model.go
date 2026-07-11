package membership

import (
	"time"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type MembershipModel struct {
	ID       membership.MembershipID
	UserID   user.UserID
	GroupID  group.GroupID
	Role     membership.Role
	JoinedAt time.Time
}

func (MembershipModel) TableName() string {
	return "memberships"
}
