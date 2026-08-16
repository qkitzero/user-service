package group

import (
	"context"
	"errors"
	"time"

	"github.com/qkitzero/user-service/internal/application/auth"
	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type GroupMembership struct {
	Group group.Group
	Role  membership.Role
}

type Member struct {
	Membership  membership.Membership
	DisplayName user.DisplayName
}

type GroupUsecase interface {
	CreateGroup(ctx context.Context, name group.GroupName) (group.Group, error)
	GetGroup(ctx context.Context, groupID group.GroupID) (group.Group, error)
	UpdateGroup(ctx context.Context, groupID group.GroupID, name group.GroupName) (group.Group, error)
	DeleteGroup(ctx context.Context, groupID group.GroupID) error
	AddChildGroup(ctx context.Context, parentID, childID group.GroupID) error
	RemoveChildGroup(ctx context.Context, parentID, childID group.GroupID) error
	ListChildGroups(ctx context.Context, groupID group.GroupID) ([]group.Group, error)
	ListParentGroups(ctx context.Context, groupID group.GroupID) ([]group.Group, error)
	AddMember(ctx context.Context, groupID group.GroupID, targetUserID user.UserID, role membership.Role) (Member, error)
	RemoveMember(ctx context.Context, groupID group.GroupID, targetUserID user.UserID) error
	UpdateMemberRole(ctx context.Context, groupID group.GroupID, targetUserID user.UserID, role membership.Role) (Member, error)
	ListMembers(ctx context.Context, groupID group.GroupID) ([]Member, error)
	ListMyGroups(ctx context.Context) ([]GroupMembership, error)
}

type groupUsecase struct {
	authService    auth.AuthService
	groupRepo      group.GroupRepository
	membershipRepo membership.MembershipRepository
	userRepo       user.UserRepository
}

func NewGroupUsecase(
	authService auth.AuthService,
	groupRepo group.GroupRepository,
	membershipRepo membership.MembershipRepository,
	userRepo user.UserRepository,
) GroupUsecase {
	return &groupUsecase{
		authService:    authService,
		groupRepo:      groupRepo,
		membershipRepo: membershipRepo,
		userRepo:       userRepo,
	}
}

func (u *groupUsecase) currentUserID(ctx context.Context) (user.UserID, error) {
	identityID, err := u.authService.VerifyToken(ctx)
	if err != nil {
		return user.UserID{}, err
	}

	id, err := identity.NewIdentityID(identityID)
	if err != nil {
		return user.UserID{}, err
	}

	caller, err := u.userRepo.FindByIdentityID(ctx, id)
	if errors.Is(err, identity.ErrIdentityNotFound) {
		return user.UserID{}, user.ErrUserNotFound
	}
	if err != nil {
		return user.UserID{}, err
	}

	return caller.ID(), nil
}

func (u *groupUsecase) toMembers(ctx context.Context, memberships []membership.Membership) ([]Member, error) {
	userIDs := make([]user.UserID, 0, len(memberships))
	for _, m := range memberships {
		userIDs = append(userIDs, m.UserID())
	}

	users, err := u.userRepo.FindByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	displayNameByUserID := make(map[user.UserID]user.DisplayName, len(users))
	for _, foundUser := range users {
		displayNameByUserID[foundUser.ID()] = foundUser.DisplayName()
	}

	members := make([]Member, 0, len(memberships))
	for _, m := range memberships {
		members = append(members, Member{Membership: m, DisplayName: displayNameByUserID[m.UserID()]})
	}

	return members, nil
}

func (u *groupUsecase) operatorRole(ctx context.Context, groupID group.GroupID) (membership.Role, error) {
	operatorID, err := u.currentUserID(ctx)
	if err != nil {
		return "", err
	}

	operator, err := u.membershipRepo.FindByUserAndGroup(ctx, operatorID, groupID)
	if err != nil {
		return "", err
	}

	return operator.Role(), nil
}

func (u *groupUsecase) CreateGroup(ctx context.Context, name group.GroupName) (group.Group, error) {
	ownerID, err := u.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	newGroup := group.NewGroup(group.NewGroupID(), name, now, now)

	if err := u.groupRepo.Create(ctx, newGroup, ownerID); err != nil {
		return nil, err
	}

	return newGroup, nil
}

func (u *groupUsecase) GetGroup(ctx context.Context, groupID group.GroupID) (group.Group, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	return u.groupRepo.FindByID(ctx, groupID)
}

func (u *groupUsecase) UpdateGroup(ctx context.Context, groupID group.GroupID, name group.GroupName) (group.Group, error) {
	role, err := u.operatorRole(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if !role.CanManageMembers() {
		return nil, membership.ErrPermissionDenied
	}

	foundGroup, err := u.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	foundGroup.Update(name)

	if err := u.groupRepo.Update(ctx, foundGroup); err != nil {
		return nil, err
	}

	return foundGroup, nil
}

func (u *groupUsecase) DeleteGroup(ctx context.Context, groupID group.GroupID) error {
	role, err := u.operatorRole(ctx, groupID)
	if err != nil {
		return err
	}
	if !role.IsOwner() {
		return membership.ErrPermissionDenied
	}

	return u.groupRepo.Delete(ctx, groupID)
}

func (u *groupUsecase) AddChildGroup(ctx context.Context, parentID, childID group.GroupID) error {
	role, err := u.operatorRole(ctx, parentID)
	if err != nil {
		return err
	}
	if !role.CanManageMembers() {
		return membership.ErrPermissionDenied
	}

	if parentID == childID {
		return group.ErrCircularReference
	}

	if _, err := u.groupRepo.FindByID(ctx, childID); err != nil {
		return err
	}

	children, err := u.groupRepo.FindChildren(ctx, parentID)
	if err != nil {
		return err
	}
	for _, child := range children {
		if child.ID() == childID {
			return group.ErrAlreadyChild
		}
	}

	descendants, err := u.collectDescendants(ctx, childID)
	if err != nil {
		return err
	}
	if descendants[parentID] {
		return group.ErrCircularReference
	}

	return u.groupRepo.AddChild(ctx, parentID, childID)
}

func (u *groupUsecase) RemoveChildGroup(ctx context.Context, parentID, childID group.GroupID) error {
	role, err := u.operatorRole(ctx, parentID)
	if err != nil {
		return err
	}
	if !role.CanManageMembers() {
		return membership.ErrPermissionDenied
	}

	return u.groupRepo.RemoveChild(ctx, parentID, childID)
}

func (u *groupUsecase) ListChildGroups(ctx context.Context, groupID group.GroupID) ([]group.Group, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	return u.groupRepo.FindChildren(ctx, groupID)
}

func (u *groupUsecase) ListParentGroups(ctx context.Context, groupID group.GroupID) ([]group.Group, error) {
	if _, err := u.authService.VerifyToken(ctx); err != nil {
		return nil, err
	}

	return u.groupRepo.FindParents(ctx, groupID)
}

func (u *groupUsecase) collectDescendants(ctx context.Context, root group.GroupID) (map[group.GroupID]bool, error) {
	visited := make(map[group.GroupID]bool)
	queue := []group.GroupID{root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		children, err := u.groupRepo.FindChildren(ctx, current)
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			if !visited[child.ID()] {
				visited[child.ID()] = true
				queue = append(queue, child.ID())
			}
		}
	}
	return visited, nil
}

func (u *groupUsecase) AddMember(ctx context.Context, groupID group.GroupID, targetUserID user.UserID, role membership.Role) (Member, error) {
	operatorID, err := u.currentUserID(ctx)
	if err != nil {
		return Member{}, err
	}

	operator, err := u.membershipRepo.FindByUserAndGroup(ctx, operatorID, groupID)
	if err != nil {
		return Member{}, err
	}
	if !operator.Role().CanManageMembers() {
		return Member{}, membership.ErrPermissionDenied
	}
	if role.IsOwner() && !operator.Role().IsOwner() {
		return Member{}, membership.ErrPermissionDenied
	}

	targetUser, err := u.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return Member{}, err
	}

	_, err = u.membershipRepo.FindByUserAndGroup(ctx, targetUserID, groupID)
	if err == nil {
		return Member{}, membership.ErrAlreadyMember
	}
	if !errors.Is(err, membership.ErrMembershipNotFound) {
		return Member{}, err
	}

	newMembership := membership.NewMembership(membership.NewMembershipID(), targetUserID, groupID, role, time.Now())
	if err := u.membershipRepo.Create(ctx, newMembership); err != nil {
		return Member{}, err
	}

	return Member{Membership: newMembership, DisplayName: targetUser.DisplayName()}, nil
}

func (u *groupUsecase) RemoveMember(ctx context.Context, groupID group.GroupID, targetUserID user.UserID) error {
	operatorID, err := u.currentUserID(ctx)
	if err != nil {
		return err
	}

	operator, err := u.membershipRepo.FindByUserAndGroup(ctx, operatorID, groupID)
	if err != nil {
		return err
	}

	if operatorID != targetUserID && !operator.Role().CanManageMembers() {
		return membership.ErrPermissionDenied
	}

	target, err := u.membershipRepo.FindByUserAndGroup(ctx, targetUserID, groupID)
	if err != nil {
		return err
	}

	if target.Role().IsOwner() {
		count, err := u.membershipRepo.CountOwners(ctx, groupID)
		if err != nil {
			return err
		}
		if count <= 1 {
			return membership.ErrLastOwnerCannotLeave
		}
	}

	return u.membershipRepo.Delete(ctx, targetUserID, groupID)
}

func (u *groupUsecase) UpdateMemberRole(ctx context.Context, groupID group.GroupID, targetUserID user.UserID, role membership.Role) (Member, error) {
	operatorID, err := u.currentUserID(ctx)
	if err != nil {
		return Member{}, err
	}

	operator, err := u.membershipRepo.FindByUserAndGroup(ctx, operatorID, groupID)
	if err != nil {
		return Member{}, err
	}
	if !operator.Role().IsOwner() {
		return Member{}, membership.ErrPermissionDenied
	}

	target, err := u.membershipRepo.FindByUserAndGroup(ctx, targetUserID, groupID)
	if err != nil {
		return Member{}, err
	}

	if target.Role().IsOwner() && !role.IsOwner() {
		count, err := u.membershipRepo.CountOwners(ctx, groupID)
		if err != nil {
			return Member{}, err
		}
		if count <= 1 {
			return Member{}, membership.ErrLastOwnerCannotLeave
		}
	}

	targetUser, err := u.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return Member{}, err
	}

	target.ChangeRole(role)
	if err := u.membershipRepo.Update(ctx, target); err != nil {
		return Member{}, err
	}

	return Member{Membership: target, DisplayName: targetUser.DisplayName()}, nil
}

func (u *groupUsecase) ListMembers(ctx context.Context, groupID group.GroupID) ([]Member, error) {
	if _, err := u.currentUserID(ctx); err != nil {
		return nil, err
	}

	memberships, err := u.membershipRepo.ListByGroupID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return u.toMembers(ctx, memberships)
}

func (u *groupUsecase) ListMyGroups(ctx context.Context) ([]GroupMembership, error) {
	userID, err := u.currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	memberships, err := u.membershipRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	groupIDs := make([]group.GroupID, 0, len(memberships))
	for _, m := range memberships {
		groupIDs = append(groupIDs, m.GroupID())
	}

	groups, err := u.groupRepo.FindByIDs(ctx, groupIDs)
	if err != nil {
		return nil, err
	}

	groupByID := make(map[group.GroupID]group.Group, len(groups))
	for _, g := range groups {
		groupByID[g.ID()] = g
	}

	groupMemberships := make([]GroupMembership, 0, len(memberships))
	for _, m := range memberships {
		g, ok := groupByID[m.GroupID()]
		if !ok {
			continue
		}
		groupMemberships = append(groupMemberships, GroupMembership{Group: g, Role: m.Role()})
	}

	return groupMemberships, nil
}
