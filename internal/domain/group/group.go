package group

import "time"

type Group interface {
	ID() GroupID
	Name() GroupName
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Update(name GroupName)
}

type group struct {
	id        GroupID
	name      GroupName
	createdAt time.Time
	updatedAt time.Time
}

func (g group) ID() GroupID {
	return g.id
}

func (g group) Name() GroupName {
	return g.name
}

func (g group) CreatedAt() time.Time {
	return g.createdAt
}

func (g group) UpdatedAt() time.Time {
	return g.updatedAt
}

func (g *group) Update(name GroupName) {
	g.name = name
	g.updatedAt = time.Now()
}

func NewGroup(
	id GroupID,
	name GroupName,
	createdAt time.Time,
	updatedAt time.Time,
) Group {
	return &group{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}
