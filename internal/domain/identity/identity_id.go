package identity

import (
	"fmt"
)

type IdentityID string

func NewIdentityID(s string) (IdentityID, error) {
	if s == "" {
		return IdentityID(""), fmt.Errorf("identity id is empty")
	}
	return IdentityID(s), nil
}
