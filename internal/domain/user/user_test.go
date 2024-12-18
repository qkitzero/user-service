package user

import "testing"

func TestNewUser(t *testing.T) {
	t.Parallel()
	email, err := NewEmail("test@example.com")
	if err != nil {
		t.Errorf("failed to new email: %v", err)
	}
	tests := []struct {
		name    string
		success bool
		id      UserID
		email   Email
	}{
		{"success new user", true, NewUserID(), email},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.id, tt.email)
			if tt.success && user == nil {
				t.Errorf("NewUser() = nil")
			}
			if tt.success && user.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", user.ID(), tt.id)
			}
			if tt.success && user.Email() != tt.email {
				t.Errorf("Email() = %v, want %v", user.Email(), tt.email)
			}
		})
	}
}
