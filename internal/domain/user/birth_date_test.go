package user

import "testing"

func TestNewBirthDate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		year    int32
		month   int32
		day     int32
	}{
		{"success new birth date", true, 2000, 1, 1},
		{"failure birth date is in the future", false, 2500, 1, 1},
		{"failure birth date is too far in the past", false, 1800, 1, 1},
		{"failure invalid birth date", false, 2000, 2, 30},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewBirthDate(tt.year, tt.month, tt.day)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error but got nil")
			}
		})
	}
}
