package identity

import (
	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type IdentityModel struct {
	ID     identity.IdentityID
	UserID user.UserID
}

func (IdentityModel) TableName() string {
	return "identities"
}
