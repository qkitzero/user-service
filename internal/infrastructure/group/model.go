package group

import (
	"time"

	"github.com/qkitzero/user-service/internal/domain/group"
)

type GroupModel struct {
	ID        group.GroupID
	Name      group.GroupName
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (GroupModel) TableName() string {
	return "groups"
}

type GroupRelationModel struct {
	ParentID group.GroupID `gorm:"primaryKey"`
	ChildID  group.GroupID `gorm:"primaryKey"`
}

func (GroupRelationModel) TableName() string {
	return "group_relations"
}
