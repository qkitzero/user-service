package membership

import "testing"

func TestNewMembershipIDFromString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		id      string
	}{
		{"success new membership id", true, "123e4567-e89b-12d3-a456-426614174000"},
		{"failure invalid uuid", false, "invalid-uuid"},
		{"failure empty", false, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			membershipID, err := NewMembershipIDFromString(tt.id)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && membershipID.String() != tt.id {
				t.Errorf("String() = %v, want %v", membershipID.String(), tt.id)
			}
		})
	}
}
