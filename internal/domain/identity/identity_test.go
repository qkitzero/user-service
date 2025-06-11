package identity

import (
	"testing"
)

func TestNewIdentity(t *testing.T) {
	t.Parallel()
	id, err := NewIdentityID("google-oauth2|000000000000000000000")
	if err != nil {
		t.Errorf("failed to new identity id: %v", err)
	}
	tests := []struct {
		name    string
		success bool
		id      IdentityID
	}{
		{"success new identity", true, id},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := NewIdentity(tt.id)
			if tt.success && identity.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", identity.ID(), tt.id)
			}
		})
	}
}
