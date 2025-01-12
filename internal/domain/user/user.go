package user

import "time"

type User interface {
	ID() UserID
	DisplayName() DisplayName
	CreatedAt() time.Time
}

type user struct {
	id          UserID
	displayName DisplayName
	createdAt   time.Time
}

func (u user) ID() UserID {
	return u.id
}

func (u user) DisplayName() DisplayName {
	return u.displayName
}

func (u user) CreatedAt() time.Time {
	return u.createdAt
}

func NewUser(id UserID, displayName DisplayName, createdAt time.Time) User {
	return user{
		id:          id,
		displayName: displayName,
		createdAt:   createdAt,
	}
}
