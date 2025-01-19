package user

import (
	"time"

	"github.com/qkitzero/user/internal/domain/user"
)

type UserTable struct {
	ID          user.UserID
	DisplayName user.DisplayName
	CreatedAt   time.Time
}

func (UserTable) TableName() string {
	return "user"
}
