package user

import "github.com/google/uuid"

type UserID struct {
	uuid.UUID
}

func NewUserID() UserID {
	return UserID{uuid.New()}
}
