package membership

import "errors"

var (
	ErrMembershipNotFound   = errors.New("membership not found")
	ErrAlreadyMember        = errors.New("already a member")
	ErrLastOwnerCannotLeave = errors.New("last owner cannot leave")
	ErrPermissionDenied     = errors.New("permission denied")
)
