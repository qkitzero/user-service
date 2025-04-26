package user

import "testing"

func TestNewUserID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
	}{
		{"success new user id", true, "fe8c2263-bbac-4bb9-a41d-b04f5afc4425"},
		{"failure empty user id", false, ""},
		{"failure invalid user id", false, "0123456789"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUserID(tt.id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
