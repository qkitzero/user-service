package identity

import "testing"

func TestNewIdentityID(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
	}{
		{"success new identity id", true, "google-oauth2|000000000000000000000"},
		{"failure empty identity id", false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewIdentityID(tt.id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
