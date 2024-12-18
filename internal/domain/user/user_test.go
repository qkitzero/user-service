package user

import "testing"

func TestNewUser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		success bool
		id      UserID
		name    string
		email   string
	}{
		{true, NewUserID(), "name", "email"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.id, tt.name, tt.email)
			if user.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", user.ID(), tt.id)
			}
			if user.Name() != tt.name {
				t.Errorf("Name() = %v, want %v", user.Name(), tt.name)
			}
			if user.Email() != tt.email {
				t.Errorf("Email() = %v, want %v", user.Email(), tt.email)
			}
		})
	}
}
