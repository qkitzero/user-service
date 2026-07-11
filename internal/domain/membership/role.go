package membership

import "fmt"

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

func (r Role) String() string {
	return string(r)
}

func NewRole(s string) (Role, error) {
	switch Role(s) {
	case RoleOwner, RoleAdmin, RoleMember:
		return Role(s), nil
	default:
		return Role(""), fmt.Errorf("invalid role: %s", s)
	}
}

func (r Role) IsOwner() bool {
	return r == RoleOwner
}

func (r Role) CanManageMembers() bool {
	return r == RoleOwner || r == RoleAdmin
}
