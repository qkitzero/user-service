package identity

type Identity interface {
	ID() IdentityID
}

type identity struct {
	id IdentityID
}

func (i identity) ID() IdentityID {
	return i.id
}

func NewIdentity(
	id IdentityID,
) Identity {
	return &identity{
		id: id,
	}
}
