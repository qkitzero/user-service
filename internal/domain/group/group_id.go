package group

import (
	"fmt"

	"github.com/google/uuid"
)

type GroupID struct {
	uuid.UUID
}

func NewGroupID() GroupID {
	id := uuid.New()
	return GroupID{id}
}

func NewGroupIDFromString(s string) (GroupID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return GroupID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return GroupID{id}, nil
}
