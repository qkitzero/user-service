package user

type User interface {
	ID() UserID
	DisplayName() DisplayName
	Email() Email
}

type user struct {
	id          UserID
	displayName DisplayName
	email       Email
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

func NewUser(id UserID, displayName DisplayName, email Email) User {
	return user{
		id:          id,
		displayName: displayName,
		email:       email,
	}
}
