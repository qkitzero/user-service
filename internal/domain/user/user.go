package user

type User interface {
	ID() UserID
	Email() Email
}

type user struct {
	id    UserID
	email Email
}

func (u user) ID() UserID {
	return u.id
}

func (u user) Email() Email {
	return u.email
}

func NewUser(id UserID, email Email) User {
	return user{
		id:    id,
		email: email,
	}
}
