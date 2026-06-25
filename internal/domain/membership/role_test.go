package membership

import "testing"

func TestNewRole(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		role    string
	}{
		{"success owner", true, "owner"},
		{"success admin", true, "admin"},
		{"success member", true, "member"},
		{"failure invalid role", false, "invalid"},
		{"failure empty", false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			role, err := NewRole(tt.role)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && role.String() != tt.role {
				t.Errorf("String() = %v, want %v", role.String(), tt.role)
			}
		})
	}
}

func TestRolePermissions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		role            Role
		isOwner         bool
		canManageMember bool
	}{
		{"owner", RoleOwner, true, true},
		{"admin", RoleAdmin, false, true},
		{"member", RoleMember, false, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.role.IsOwner() != tt.isOwner {
				t.Errorf("IsOwner() = %v, want %v", tt.role.IsOwner(), tt.isOwner)
			}
			if tt.role.CanManageMembers() != tt.canManageMember {
				t.Errorf("CanManageMembers() = %v, want %v", tt.role.CanManageMembers(), tt.canManageMember)
			}
		})
	}
}
