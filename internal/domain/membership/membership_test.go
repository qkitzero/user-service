package membership

import (
	"testing"
	"time"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/user"
)

func TestNewMembership(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		role Role
	}{
		{"owner", RoleOwner},
		{"admin", RoleAdmin},
		{"member", RoleMember},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewMembershipID()
			userID := user.NewUserID()
			groupID := group.NewGroupID()
			now := time.Now()

			m := NewMembership(id, userID, groupID, tt.role, now)

			if m.ID() != id {
				t.Errorf("ID() = %v, want %v", m.ID(), id)
			}
			if m.UserID() != userID {
				t.Errorf("UserID() = %v, want %v", m.UserID(), userID)
			}
			if m.GroupID() != groupID {
				t.Errorf("GroupID() = %v, want %v", m.GroupID(), groupID)
			}
			if m.Role() != tt.role {
				t.Errorf("Role() = %v, want %v", m.Role(), tt.role)
			}
			if !m.JoinedAt().Equal(now) {
				t.Errorf("JoinedAt() = %v, want %v", m.JoinedAt(), now)
			}
		})
	}
}

func TestMembershipChangeRole(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		initial Role
		changed Role
	}{
		{"member to admin", RoleMember, RoleAdmin},
		{"admin to owner", RoleAdmin, RoleOwner},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := NewMembership(NewMembershipID(), user.NewUserID(), group.NewGroupID(), tt.initial, time.Now())

			m.ChangeRole(tt.changed)

			if m.Role() != tt.changed {
				t.Errorf("Role() = %v, want %v", m.Role(), tt.changed)
			}
		})
	}
}
