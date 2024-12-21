package user

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	displayName, err := NewDisplayName("test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	email, err := NewEmail("test@example.com")
	if err != nil {
		t.Errorf("failed to new email: %v", err)
	}
	tests := []struct {
		name        string
		success     bool
		id          UserID
		displayName DisplayName
		email       Email
	}{
		{"success new user", true, NewUserID(), displayName, email},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.id, tt.displayName, tt.email)
			if tt.success && user == nil {
				t.Errorf("NewUser() = nil")
			}
			if tt.success && user.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", user.ID(), tt.id)
			}
			if tt.success && user.DisplayName() != tt.displayName {
				t.Errorf("DisplayName() = %v, want %v", user.DisplayName(), tt.displayName)
			}
			if tt.success && user.Email() != tt.email {
				t.Errorf("Email() = %v, want %v", user.Email(), tt.email)
			}
			if tt.success && time.Now().Sub(user.CreatedAt()) > time.Second {
				t.Errorf("CreatedAt() = %v, want close to current time", user.CreatedAt())
			}
		})
	}
}
