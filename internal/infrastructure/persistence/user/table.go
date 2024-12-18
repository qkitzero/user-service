package user

import "user/internal/domain/user"

type UserTable struct {
	ID    user.UserID
	Name  string
	Email string
}

func (UserTable) TableName() string {
	return "user"
}
