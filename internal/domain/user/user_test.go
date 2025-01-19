package user

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	id, err := NewUserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	displayName, err := NewDisplayName("test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	tests := []struct {
		name        string
		success     bool
		id          UserID
		displayName DisplayName
		createdAt   time.Time
	}{
		{"success new user", true, id, displayName, time.Now()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.id, tt.displayName, tt.createdAt)
			if tt.success && user == nil {
				t.Errorf("NewUser() = nil")
			}
			if tt.success && user.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", user.ID(), tt.id)
			}
			if tt.success && user.DisplayName() != tt.displayName {
				t.Errorf("DisplayName() = %v, want %v", user.DisplayName(), tt.displayName)
			}
			if tt.success && time.Now().Sub(user.CreatedAt()) > time.Second {
				t.Errorf("CreatedAt() = %v, want close to current time", user.CreatedAt())
			}
		})
	}
}
