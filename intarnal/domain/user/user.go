package user

type User interface {
	ID() string
	Name() string
	Email() string
}

type user struct {
	id    string
	name  string
	email string
}

func (u user) ID() string {
	return u.id
}

func (u user) Name() string {
	return u.name
}

func (u user) Email() string {
	return u.email
}

func NewUser(id, name, email string) User {
	return user{
		id:    id,
		name:  name,
		email: email,
	}
}
