package user

import (
	"time"

	"github.com/qkitzero/user-service/internal/domain/user"
)

type UserModel struct {
	ID          user.UserID
	DisplayName user.DisplayName
	BirthDate   user.BirthDate
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (UserModel) TableName() string {
	return "users"
}
