package user

import "testing"

func TestNewEmail(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		email   string
	}{
		{"success new email", true, "test@example.com"},
		{"failure invalid email", false, "invalid-email"},
		{"failure empty email", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.email)
			if tt.success && err != nil {
				t.Errorf("expected no error")
			}
			if !tt.success && err == nil {
				t.Errorf("expected error")
			}
		})
	}
}
