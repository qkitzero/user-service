package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/qkitzero/user-service/internal/domain/identity"
)

func TestNewUser(t *testing.T) {
	t.Parallel()
	id, err := NewUserIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	identityID, err := identity.NewIdentityID("google-oauth2|000000000000000000000")
	if err != nil {
		t.Errorf("failed to new identity id: %v", err)
	}
	identities := []identity.Identity{
		identity.NewIdentity(identityID),
	}
	displayName, err := NewDisplayName("test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	birthDate, err := NewBirthDate(2000, 1, 1)
	if err != nil {
		t.Errorf("failed to new birth date: %v", err)
	}
	tests := []struct {
		name        string
		success     bool
		id          UserID
		identities  []identity.Identity
		displayName DisplayName
		birthDate   BirthDate
		createdAt   time.Time
		updatedAt   time.Time
	}{
		{"success new user", true, id, identities, displayName, birthDate, time.Now(), time.Now()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := NewUser(tt.id, tt.identities, tt.displayName, tt.birthDate, tt.createdAt, tt.updatedAt)
			if tt.success && user.ID() != tt.id {
				t.Errorf("ID() = %v, want %v", user.ID(), tt.id)
			}
			if tt.success && !reflect.DeepEqual(user.Identities(), tt.identities) {
				t.Errorf("Identities() = %v, want %v", user.Identities(), tt.identities)
			}
			if tt.success && user.DisplayName() != tt.displayName {
				t.Errorf("DisplayName() = %v, want %v", user.DisplayName(), tt.displayName)
			}
			if tt.success && user.BirthDate() != tt.birthDate {
				t.Errorf("BirthDate() = %v, want %v", user.BirthDate(), tt.birthDate)
			}
			if tt.success && !user.CreatedAt().Equal(tt.createdAt) {
				t.Errorf("CreatedAt() = %v, want %v", user.CreatedAt(), tt.createdAt)
			}
			if tt.success && !user.UpdatedAt().Equal(tt.updatedAt) {
				t.Errorf("UpdatedAt() = %v, want %v", user.UpdatedAt(), tt.updatedAt)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	id, err := NewUserIDFromString("fe8c2263-bbac-4bb9-a41d-b04f5afc4425")
	if err != nil {
		t.Errorf("failed to new user id: %v", err)
	}
	identityID, err := identity.NewIdentityID("google-oauth2|000000000000000000000")
	if err != nil {
		t.Errorf("failed to new identity id: %v", err)
	}
	identities := []identity.Identity{
		identity.NewIdentity(identityID),
	}
	displayName, err := NewDisplayName("test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	birthDate, err := NewBirthDate(2000, 1, 1)
	if err != nil {
		t.Errorf("failed to new birth date: %v", err)
	}
	updatedDisplayName, err := NewDisplayName("updated test user")
	if err != nil {
		t.Errorf("failed to new display name: %v", err)
	}
	user := NewUser(id, identities, displayName, birthDate, time.Now(), time.Now())
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
			if tt.success && tt.user.DisplayName() != tt.updatedDisplayName {
				t.Errorf("DisplayName() = %v, want %v", tt.user.DisplayName(), tt.updatedDisplayName)
			}
			if tt.success && !tt.user.CreatedAt().Before(tt.user.UpdatedAt()) {
				t.Errorf("CreatedAt() = %v, UpdatedAt() = %v, want CreatedAt < UpdatedAt", tt.user.CreatedAt(), tt.user.UpdatedAt())
			}
		})
	}
}
