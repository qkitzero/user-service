package membership

import (
	"fmt"

	"github.com/google/uuid"
)

type MembershipID struct {
	uuid.UUID
}

func NewMembershipID() MembershipID {
	id := uuid.New()
	return MembershipID{id}
}

func NewMembershipIDFromString(s string) (MembershipID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return MembershipID{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return MembershipID{id}, nil
}
