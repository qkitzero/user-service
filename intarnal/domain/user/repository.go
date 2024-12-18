package user

type UserRepository interface {
	Create(user User) error
}
