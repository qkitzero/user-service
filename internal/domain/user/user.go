package user

import "time"

type User interface {
	ID() UserID
	DisplayName() DisplayName
	Email() Email
	CreatedAt() time.Time
}

type user struct {
	id          UserID
	displayName DisplayName
	email       Email
	createdAt   time.Time
}

func (u user) ID() UserID {
	return u.id
}

func (u user) DisplayName() DisplayName {
	return u.displayName
}

func (u user) Email() Email {
	return u.email
}

func (u user) CreatedAt() time.Time {
	return u.createdAt
}

func NewUser(id UserID, displayName DisplayName, email Email) User {
	now := time.Now()
	return user{
		id:          id,
		displayName: displayName,
		email:       email,
		createdAt:   now,
	}
}
