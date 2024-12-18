package user

type User interface {
	ID() UserID
	Name() string
	Email() string
}

type user struct {
	id    UserID
	name  string
	email string
}

func (u user) ID() UserID {
	return u.id
}

func (u user) Name() string {
	return u.name
}

func (u user) Email() string {
	return u.email
}

func NewUser(id UserID, name, email string) User {
	return user{
		id:    id,
		name:  name,
		email: email,
	}
}
