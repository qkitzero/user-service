package group

import (
	"testing"
	"time"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		groupName string
	}{
		{"success new group", "test group"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id := NewGroupID()
			groupName, _ := NewGroupName(tt.groupName)
			now := time.Now()

			g := NewGroup(id, groupName, now, now)

			if g.ID() != id {
				t.Errorf("ID() = %v, want %v", g.ID(), id)
			}
			if g.Name() != groupName {
				t.Errorf("Name() = %v, want %v", g.Name(), groupName)
			}
			if !g.CreatedAt().Equal(now) {
				t.Errorf("CreatedAt() = %v, want %v", g.CreatedAt(), now)
			}
			if !g.UpdatedAt().Equal(now) {
				t.Errorf("UpdatedAt() = %v, want %v", g.UpdatedAt(), now)
			}
		})
	}
}

func TestGroupUpdate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		newName string
	}{
		{"success update name", "updated group"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, _ := NewGroupName("test group")
			now := time.Now()
			g := NewGroup(NewGroupID(), name, now, now)

			newName, _ := NewGroupName(tt.newName)
			g.Update(newName)

			if g.Name() != newName {
				t.Errorf("Name() = %v, want %v", g.Name(), newName)
			}
		})
	}
}
