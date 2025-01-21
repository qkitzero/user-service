package user

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	id, err := NewUserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	displayName, err := NewDisplayName("test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	tests := []struct {
		name        string
		success     bool
		id          UserID
		displayName DisplayName
		createdAt   time.Time
		updatedAt   time.Time
	}{
		{"success new user", true, id, displayName, time.Now(), time.Now()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.id, tt.displayName, tt.createdAt, tt.updatedAt)
			if tt.success && user == nil {
				t.Errorf("NewUser() = nil")
			}
			if tt.success && user.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", user.ID(), tt.id)
			}
			if tt.success && user.DisplayName() != tt.displayName {
				t.Errorf("DisplayName() = %v, want %v", user.DisplayName(), tt.displayName)
			}
			if tt.success && time.Now().Sub(user.CreatedAt()) > time.Second {
				t.Errorf("CreatedAt() = %v, want close to current time", user.CreatedAt())
			}
			if tt.success && time.Now().Sub(user.UpdatedAt()) > time.Second {
				t.Errorf("UpdatedAt() = %v, want close to current time", user.UpdatedAt())
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	id, err := NewUserID("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	displayName, err := NewDisplayName("test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	updatedDisplayName, err := NewDisplayName("updated test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	user := NewUser(id, displayName, time.Now(), time.Now())
	tests := []struct {
		name               string
		success            bool
		user               User
		updatedDisplayName DisplayName
	}{
		{"success update user", true, user, updatedDisplayName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.user.Update(tt.updatedDisplayName)
			if tt.success && user.DisplayName() != tt.updatedDisplayName {
				t.Errorf("DisplayName() = %v, want %v", user.DisplayName(), tt.updatedDisplayName)
			}
			if tt.success && !user.CreatedAt().Before(user.UpdatedAt()) {
				t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want CreatedAt < UpdatedAt", user.CreatedAt(), user.UpdatedAt())
			}
		})
	}
}
