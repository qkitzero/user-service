package user

import (
	"time"

	"github.com/qkitzero/user-service/internal/domain/identity"
)

type User interface {
	ID() UserID
	Identities() []identity.Identity
	DisplayName() DisplayName
	BirthDate() BirthDate
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(displayName DisplayName, birthDate BirthDate)
}

type user struct {
	id          UserID
	identities  []identity.Identity
	displayName DisplayName
	birthDate   BirthDate
	createdAt   time.Time
	updatedAt   time.Time
}

func (u user) ID() UserID {
	return u.id
}

func (u user) Identities() []identity.Identity {
	return u.identities
}

func (u user) DisplayName() DisplayName {
	return u.displayName
}

func (u user) BirthDate() BirthDate {
	return u.birthDate
}

func (u user) CreatedAt() time.Time {
	return u.createdAt
}

func (u user) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *user) Update(displayName DisplayName, birthDate BirthDate) {
	u.displayName = displayName
	u.birthDate = birthDate
	u.updatedAt = time.Now()
}

func NewUser(
	id UserID,
	identities []identity.Identity,
	displayName DisplayName,
	birthDate BirthDate,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return &user{
		id:          id,
		identities:  identities,
		displayName: displayName,
		birthDate:   birthDate,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
