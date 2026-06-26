package group

import (
	"strings"
	"testing"
)

func TestNewGroupName(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		success   bool
		groupName string
		want      string
	}{
		{"success new group name", true, "test group", "test group"},
		{"success trim spaces", true, "  test group  ", "test group"},
		{"failure empty group name", false, "", ""},
		{"failure whitespace only", false, "   ", ""},
		{"failure too long", false, strings.Repeat("a", maxGroupNameLength+1), ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			groupName, err := NewGroupName(tt.groupName)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && groupName.String() != tt.want {
				t.Errorf("String() = %v, want %v", groupName.String(), tt.want)
			}
		})
	}
}
