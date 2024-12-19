package user

import "user/internal/domain/user"

type UserTable struct {
	ID          user.UserID
	DisplayName user.DisplayName
	Email       user.Email
}

func (UserTable) TableName() string {
	return "user"
}
