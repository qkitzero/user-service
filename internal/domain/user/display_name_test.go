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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			displayName, err := NewDisplayName(tt.displayName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && displayName.String() != tt.displayName {
				t.Errorf("String() = %v, want %v", displayName.String(), tt.displayName)
			}
		})
	}
}
