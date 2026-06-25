package membership

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) membership.MembershipRepository {
	return &membershipRepository{db: db}
}

func toMembership(m MembershipModel) membership.Membership {
	return membership.NewMembership(m.ID, m.UserID, m.GroupID, m.Role, m.JoinedAt)
}

func (r *membershipRepository) Create(ctx context.Context, m membership.Membership) error {
	membershipModel := MembershipModel{
		ID:       m.ID(),
		UserID:   m.UserID(),
		GroupID:  m.GroupID(),
		Role:     m.Role(),
		JoinedAt: m.JoinedAt(),
	}

	return r.db.WithContext(ctx).Create(&membershipModel).Error
}

func (r *membershipRepository) FindByUserAndGroup(ctx context.Context, userID user.UserID, groupID group.GroupID) (membership.Membership, error) {
	var membershipModel MembershipModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND group_id = ?", userID, groupID).First(&membershipModel).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, membership.ErrMembershipNotFound
	}
	if err != nil {
		return nil, err
	}

	return toMembership(membershipModel), nil
}

func (r *membershipRepository) ListByGroupID(ctx context.Context, groupID group.GroupID) ([]membership.Membership, error) {
	var membershipModels []MembershipModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&membershipModels).Error; err != nil {
		return nil, err
	}

	memberships := make([]membership.Membership, 0, len(membershipModels))
	for _, m := range membershipModels {
		memberships = append(memberships, toMembership(m))
	}

	return memberships, nil
}

func (r *membershipRepository) ListByUserID(ctx context.Context, userID user.UserID) ([]membership.Membership, error) {
	var membershipModels []MembershipModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&membershipModels).Error; err != nil {
		return nil, err
	}

	memberships := make([]membership.Membership, 0, len(membershipModels))
	for _, m := range membershipModels {
		memberships = append(memberships, toMembership(m))
	}

	return memberships, nil
}

func (r *membershipRepository) Update(ctx context.Context, m membership.Membership) error {
	membershipModel := MembershipModel{
		ID:       m.ID(),
		UserID:   m.UserID(),
		GroupID:  m.GroupID(),
		Role:     m.Role(),
		JoinedAt: m.JoinedAt(),
	}

	return r.db.WithContext(ctx).Save(&membershipModel).Error
}

func (r *membershipRepository) Delete(ctx context.Context, userID user.UserID, groupID group.GroupID) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&MembershipModel{}).Error
}

func (r *membershipRepository) CountOwners(ctx context.Context, groupID group.GroupID) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&MembershipModel{}).Where("group_id = ? AND role = ?", groupID, membership.RoleOwner).Count(&count).Error
	return int(count), err
}
