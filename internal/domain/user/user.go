package user

import "time"

type User interface {
	ID() UserID
	DisplayName() DisplayName
	BirthDate() BirthDate
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(displayName DisplayName)
}

type user struct {
	id          UserID
	displayName DisplayName
	birthDate   BirthDate
	createdAt   time.Time
	updatedAt   time.Time
}

func (u user) ID() UserID {
	return u.id
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

func (u *user) Update(displayName DisplayName) {
	u.displayName = displayName
	u.updatedAt = time.Now()
}

func NewUser(
	id UserID,
	displayName DisplayName,
	birthDate BirthDate,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return &user{
		id:          id,
		displayName: displayName,
		birthDate:   birthDate,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
