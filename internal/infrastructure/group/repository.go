package group

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
	inframembership "github.com/qkitzero/user-service/internal/infrastructure/membership"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) group.GroupRepository {
	return &groupRepository{db: db}
}

func toGroup(m GroupModel) group.Group {
	return group.NewGroup(m.ID, m.Name, m.CreatedAt, m.UpdatedAt)
}

func (r *groupRepository) Create(ctx context.Context, g group.Group, ownerUserID user.UserID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		groupModel := GroupModel{
			ID:        g.ID(),
			Name:      g.Name(),
			CreatedAt: g.CreatedAt(),
			UpdatedAt: g.UpdatedAt(),
		}

		if err := tx.Create(&groupModel).Error; err != nil {
			return err
		}

		ownerModel := inframembership.MembershipModel{
			ID:       membership.NewMembershipID(),
			UserID:   ownerUserID,
			GroupID:  g.ID(),
			Role:     membership.RoleOwner,
			JoinedAt: g.CreatedAt(),
		}

		return tx.Create(&ownerModel).Error
	})
}

func (r *groupRepository) FindByID(ctx context.Context, id group.GroupID) (group.Group, error) {
	var groupModel GroupModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&groupModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, group.ErrGroupNotFound
	}
	if err != nil {
		return nil, err
	}

	return toGroup(groupModel), nil
}

func (r *groupRepository) Update(ctx context.Context, g group.Group) error {
	groupModel := GroupModel{
		ID:        g.ID(),
		Name:      g.Name(),
		CreatedAt: g.CreatedAt(),
		UpdatedAt: g.UpdatedAt(),
	}

	return r.db.WithContext(ctx).Save(&groupModel).Error
}

func (r *groupRepository) Delete(ctx context.Context, id group.GroupID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&GroupModel{}).Error
}

func (r *groupRepository) AddChild(ctx context.Context, parentID, childID group.GroupID) error {
	relationModel := GroupRelationModel{
		ParentID: parentID,
		ChildID:  childID,
	}

	return r.db.WithContext(ctx).Create(&relationModel).Error
}

func (r *groupRepository) RemoveChild(ctx context.Context, parentID, childID group.GroupID) error {
	return r.db.WithContext(ctx).Where("parent_id = ? AND child_id = ?", parentID, childID).Delete(&GroupRelationModel{}).Error
}

func (r *groupRepository) FindChildren(ctx context.Context, parentID group.GroupID) ([]group.Group, error) {
	var relations []GroupRelationModel
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&relations).Error; err != nil {
		return nil, err
	}

	childIDs := make([]group.GroupID, 0, len(relations))
	for _, rel := range relations {
		childIDs = append(childIDs, rel.ChildID)
	}

	return r.FindByIDs(ctx, childIDs)
}

func (r *groupRepository) FindParents(ctx context.Context, childID group.GroupID) ([]group.Group, error) {
	var relations []GroupRelationModel
	if err := r.db.WithContext(ctx).Where("child_id = ?", childID).Find(&relations).Error; err != nil {
		return nil, err
	}

	parentIDs := make([]group.GroupID, 0, len(relations))
	for _, rel := range relations {
		parentIDs = append(parentIDs, rel.ParentID)
	}

	return r.FindByIDs(ctx, parentIDs)
}

func (r *groupRepository) FindByIDs(ctx context.Context, ids []group.GroupID) ([]group.Group, error) {
	if len(ids) == 0 {
		return []group.Group{}, nil
	}

	var groupModels []GroupModel
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&groupModels).Error; err != nil {
		return nil, err
	}

	groups := make([]group.Group, 0, len(groupModels))
	for _, m := range groupModels {
		groups = append(groups, toGroup(m))
	}

	return groups, nil
}
