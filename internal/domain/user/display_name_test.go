package user

import "testing"

func TestNewDisplayName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		success     bool
		displayName string
	}{
		{"success new display name", true, "test user"},
		{"failure empty display name", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDisplayName(tt.displayName)
			if tt.success && err != nil {
				t.Errorf("expected no error")
			}
			if !tt.success && err == nil {
				t.Errorf("expected error")
			}
		})
	}
}
