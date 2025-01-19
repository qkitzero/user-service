package user

type UserRepository interface {
	Create(user User) error
	Read(userID UserID) (User, error)
}
