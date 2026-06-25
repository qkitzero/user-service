package membership

import (
	"time"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type Membership interface {
	ID() MembershipID
	UserID() user.UserID
	GroupID() group.GroupID
	Role() Role
	JoinedAt() time.Time
	ChangeRole(role Role)
}

type membership struct {
	id       MembershipID
	userID   user.UserID
	groupID  group.GroupID
	role     Role
	joinedAt time.Time
}

func (m membership) ID() MembershipID {
	return m.id
}

func (m membership) UserID() user.UserID {
	return m.userID
}

func (m membership) GroupID() group.GroupID {
	return m.groupID
}

func (m membership) Role() Role {
	return m.role
}

func (m membership) JoinedAt() time.Time {
	return m.joinedAt
}

func (m *membership) ChangeRole(role Role) {
	m.role = role
}

func NewMembership(
	id MembershipID,
	userID user.UserID,
	groupID group.GroupID,
	role Role,
	joinedAt time.Time,
) Membership {
	return &membership{
		id:       id,
		userID:   userID,
		groupID:  groupID,
		role:     role,
		joinedAt: joinedAt,
	}
}
