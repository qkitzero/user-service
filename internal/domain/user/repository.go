package user

import "github.com/qkitzero/user/internal/domain/identity"

type UserRepository interface {
	Create(user User) error
	FindByIdentityID(identityID identity.IdentityID) (User, error)
	FindByID(userID UserID) (User, error)
	Update(user User) error
}
