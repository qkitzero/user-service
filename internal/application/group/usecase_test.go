package group

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/identity"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksappauth "github.com/qkitzero/user-service/mocks/application/auth"
	mocksgroup "github.com/qkitzero/user-service/mocks/domain/group"
	mocksmembership "github.com/qkitzero/user-service/mocks/domain/membership"
	mocksuser "github.com/qkitzero/user-service/mocks/domain/user"
)

const validIdentityID = "google-oauth2|000000000000000000000"

func setupAuth(ctrl *gomock.Controller, mockAuth *mocksappauth.MockAuthService, mockUserRepo *mocksuser.MockUserRepository, callerID user.UserID, identityID string, verifyErr, findUserErr error) {
	mockUser := mocksuser.NewMockUser(ctrl)
	mockUser.EXPECT().ID().Return(callerID).AnyTimes()
	mockAuth.EXPECT().VerifyToken(gomock.Any()).Return(identityID, verifyErr).AnyTimes()
	mockUserRepo.EXPECT().FindByIdentityID(gomock.Any(), gomock.Eq(identity.IdentityID(identityID))).Return(mockUser, findUserErr).AnyTimes()
}

func TestCreateGroup(t *testing.T) {
	t.Parallel()
	name, _ := group.NewGroupName("test group")

	tests := []struct {
		name           string
		success        bool
		identityID     string
		verifyTokenErr error
		findUserErr    error
		createErr      error
	}{
		{"success create group", true, validIdentityID, nil, nil, nil},
		{"failure verify token error", false, validIdentityID, fmt.Errorf("verify token error"), nil, nil},
		{"failure invalid identity id", false, "", nil, nil, nil},
		{"failure find user error", false, validIdentityID, nil, errors.New("find user error"), nil},
		{"failure identity not found", false, validIdentityID, nil, identity.ErrIdentityNotFound, nil},
		{"failure create error", false, validIdentityID, nil, nil, errors.New("create error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, user.NewUserID(), tt.identityID, tt.verifyTokenErr, tt.findUserErr)
			mockGroupRepo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.CreateGroup(context.Background(), name)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestGetGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		verifyTokenErr error
		findByIDErr    error
	}{
		{"success get group", true, nil, nil},
		{"failure verify token error", false, fmt.Errorf("verify token error"), nil},
		{"failure group not found", false, nil, group.ErrGroupNotFound},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			mockGroup := mocksgroup.NewMockGroup(ctrl)
			mockAuth.EXPECT().VerifyToken(gomock.Any()).Return(validIdentityID, tt.verifyTokenErr).AnyTimes()
			mockGroupRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockGroup, tt.findByIDErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.GetGroup(context.Background(), group.NewGroupID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestUpdateGroup(t *testing.T) {
	t.Parallel()
	name, _ := group.NewGroupName("updated group")
	operatorID := user.NewUserID()

	tests := []struct {
		name         string
		success      bool
		verifyErr    error
		operatorFind error
		operatorRole membership.Role
		findByIDErr  error
		updateErr    error
	}{
		{"success update group", true, nil, nil, membership.RoleAdmin, nil, nil},
		{"failure verify token error", false, errors.New("verify"), nil, membership.RoleAdmin, nil, nil},
		{"failure not member", false, nil, membership.ErrMembershipNotFound, membership.RoleAdmin, nil, nil},
		{"failure permission denied", false, nil, nil, membership.RoleMember, nil, nil},
		{"failure group not found", false, nil, nil, membership.RoleAdmin, group.ErrGroupNotFound, nil},
		{"failure update error", false, nil, nil, membership.RoleAdmin, nil, errors.New("update")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, operatorID, validIdentityID, tt.verifyErr, nil)

			operatorMembership := mocksmembership.NewMockMembership(ctrl)
			operatorMembership.EXPECT().Role().Return(tt.operatorRole).AnyTimes()
			mockMembershipRepo.EXPECT().FindByUserAndGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(operatorMembership, tt.operatorFind).AnyTimes()

			mockGroup := mocksgroup.NewMockGroup(ctrl)
			mockGroup.EXPECT().Update(gomock.Any()).AnyTimes()
			mockGroupRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(mockGroup, tt.findByIDErr).AnyTimes()
			mockGroupRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.UpdateGroup(context.Background(), group.NewGroupID(), name)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestDeleteGroup(t *testing.T) {
	t.Parallel()
	operatorID := user.NewUserID()

	tests := []struct {
		name         string
		success      bool
		verifyErr    error
		operatorFind error
		operatorRole membership.Role
		deleteErr    error
	}{
		{"success delete group", true, nil, nil, membership.RoleOwner, nil},
		{"failure verify token error", false, errors.New("verify"), nil, membership.RoleOwner, nil},
		{"failure not member", false, nil, membership.ErrMembershipNotFound, membership.RoleOwner, nil},
		{"failure permission denied", false, nil, nil, membership.RoleAdmin, nil},
		{"failure delete error", false, nil, nil, membership.RoleOwner, errors.New("delete")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, operatorID, validIdentityID, tt.verifyErr, nil)

			operatorMembership := mocksmembership.NewMockMembership(ctrl)
			operatorMembership.EXPECT().Role().Return(tt.operatorRole).AnyTimes()
			mockMembershipRepo.EXPECT().FindByUserAndGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(operatorMembership, tt.operatorFind).AnyTimes()
			mockGroupRepo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(tt.deleteErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			err := u.DeleteGroup(context.Background(), group.NewGroupID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestAddChildGroup(t *testing.T) {
	t.Parallel()
	operatorID := user.NewUserID()

	tests := []struct {
		name            string
		success         bool
		verifyErr       error
		operatorFind    error
		operatorRole    membership.Role
		sameID          bool
		childFindErr    error
		findChildrenErr error
		alreadyChild    bool
		descendantsErr  error
		circular        bool
		addChildErr     error
	}{
		{"success add child", true, nil, nil, membership.RoleAdmin, false, nil, nil, false, nil, false, nil},
		{"failure verify token error", false, errors.New("verify"), nil, membership.RoleAdmin, false, nil, nil, false, nil, false, nil},
		{"failure not member", false, nil, membership.ErrMembershipNotFound, membership.RoleAdmin, false, nil, nil, false, nil, false, nil},
		{"failure permission denied", false, nil, nil, membership.RoleMember, false, nil, nil, false, nil, false, nil},
		{"failure same id", false, nil, nil, membership.RoleAdmin, true, nil, nil, false, nil, false, nil},
		{"failure child not found", false, nil, nil, membership.RoleAdmin, false, group.ErrGroupNotFound, nil, false, nil, false, nil},
		{"failure find children error", false, nil, nil, membership.RoleAdmin, false, nil, errors.New("find children error"), false, nil, false, nil},
		{"failure already child", false, nil, nil, membership.RoleAdmin, false, nil, nil, true, nil, false, nil},
		{"failure collect descendants error", false, nil, nil, membership.RoleAdmin, false, nil, nil, false, errors.New("collect descendants error"), false, nil},
		{"failure circular reference", false, nil, nil, membership.RoleAdmin, false, nil, nil, false, nil, true, nil},
		{"failure add child error", false, nil, nil, membership.RoleAdmin, false, nil, nil, false, nil, false, errors.New("add child error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			parentID := group.NewGroupID()
			childID := group.NewGroupID()
			if tt.sameID {
				childID = parentID
			}

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, operatorID, validIdentityID, tt.verifyErr, nil)

			operatorMembership := mocksmembership.NewMockMembership(ctrl)
			operatorMembership.EXPECT().Role().Return(tt.operatorRole).AnyTimes()
			mockMembershipRepo.EXPECT().FindByUserAndGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(operatorMembership, tt.operatorFind).AnyTimes()

			childGroup := mocksgroup.NewMockGroup(ctrl)
			mockGroupRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(childGroup, tt.childFindErr).AnyTimes()

			mockGroupRepo.EXPECT().FindChildren(gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, gid group.GroupID) ([]group.Group, error) {
					if tt.findChildrenErr != nil && gid == parentID {
						return nil, tt.findChildrenErr
					}
					if tt.descendantsErr != nil && gid == childID {
						return nil, tt.descendantsErr
					}
					if tt.alreadyChild && gid == parentID {
						c := mocksgroup.NewMockGroup(ctrl)
						c.EXPECT().ID().Return(childID).AnyTimes()
						return []group.Group{c}, nil
					}
					if tt.circular && gid == childID {
						p := mocksgroup.NewMockGroup(ctrl)
						p.EXPECT().ID().Return(parentID).AnyTimes()
						return []group.Group{p}, nil
					}
					return []group.Group{}, nil
				}).AnyTimes()
			mockGroupRepo.EXPECT().AddChild(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.addChildErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			err := u.AddChildGroup(context.Background(), parentID, childID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestRemoveChildGroup(t *testing.T) {
	t.Parallel()
	operatorID := user.NewUserID()

	tests := []struct {
		name           string
		success        bool
		verifyErr      error
		operatorFind   error
		operatorRole   membership.Role
		removeChildErr error
	}{
		{"success remove child", true, nil, nil, membership.RoleAdmin, nil},
		{"failure verify token error", false, errors.New("verify"), nil, membership.RoleAdmin, nil},
		{"failure not member", false, nil, membership.ErrMembershipNotFound, membership.RoleAdmin, nil},
		{"failure permission denied", false, nil, nil, membership.RoleMember, nil},
		{"failure remove child error", false, nil, nil, membership.RoleAdmin, errors.New("remove child error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, operatorID, validIdentityID, tt.verifyErr, nil)

			operatorMembership := mocksmembership.NewMockMembership(ctrl)
			operatorMembership.EXPECT().Role().Return(tt.operatorRole).AnyTimes()
			mockMembershipRepo.EXPECT().FindByUserAndGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(operatorMembership, tt.operatorFind).AnyTimes()
			mockGroupRepo.EXPECT().RemoveChild(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.removeChildErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			err := u.RemoveChildGroup(context.Background(), group.NewGroupID(), group.NewGroupID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListChildGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		success         bool
		verifyTokenErr  error
		findChildrenErr error
	}{
		{"success list child groups", true, nil, nil},
		{"failure verify token error", false, fmt.Errorf("verify token error"), nil},
		{"failure find children error", false, nil, errors.New("find children error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			mockAuth.EXPECT().VerifyToken(gomock.Any()).Return(validIdentityID, tt.verifyTokenErr).AnyTimes()
			mockGroupRepo.EXPECT().FindChildren(gomock.Any(), gomock.Any()).Return(nil, tt.findChildrenErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.ListChildGroups(context.Background(), group.NewGroupID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListParentGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		success        bool
		verifyTokenErr error
		findParentsErr error
	}{
		{"success list parent groups", true, nil, nil},
		{"failure verify token error", false, fmt.Errorf("verify token error"), nil},
		{"failure find parents error", false, nil, errors.New("find parents error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			mockAuth.EXPECT().VerifyToken(gomock.Any()).Return(validIdentityID, tt.verifyTokenErr).AnyTimes()
			mockGroupRepo.EXPECT().FindParents(gomock.Any(), gomock.Any()).Return(nil, tt.findParentsErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.ListParentGroups(context.Background(), group.NewGroupID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestAddMember(t *testing.T) {
	t.Parallel()
	operatorID := user.NewUserID()
	targetID := user.NewUserID()
	groupID := group.NewGroupID()

	tests := []struct {
		name          string
		success       bool
		verifyErr     error
		operatorRole  membership.Role
		operatorFind  error
		grantRole     membership.Role
		targetUserErr error
		targetExists  bool
		createErr     error
	}{
		{"success add member", true, nil, membership.RoleAdmin, nil, membership.RoleMember, nil, false, nil},
		{"failure verify token", false, errors.New("verify"), membership.RoleAdmin, nil, membership.RoleMember, nil, false, nil},
		{"failure operator not member", false, nil, membership.RoleAdmin, membership.ErrMembershipNotFound, membership.RoleMember, nil, false, nil},
		{"failure permission denied", false, nil, membership.RoleMember, nil, membership.RoleMember, nil, false, nil},
		{"failure owner grant by admin", false, nil, membership.RoleAdmin, nil, membership.RoleOwner, nil, false, nil},
		{"failure target user not found", false, nil, membership.RoleAdmin, nil, membership.RoleMember, user.ErrUserNotFound, false, nil},
		{"failure already member", false, nil, membership.RoleAdmin, nil, membership.RoleMember, nil, true, nil},
		{"failure create error", false, nil, membership.RoleAdmin, nil, membership.RoleMember, nil, false, errors.New("create")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, operatorID, validIdentityID, tt.verifyErr, nil)

			targetUser := mocksuser.NewMockUser(ctrl)
			mockUserRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(targetUser, tt.targetUserErr).AnyTimes()

			operatorMembership := mocksmembership.NewMockMembership(ctrl)
			operatorMembership.EXPECT().Role().Return(tt.operatorRole).AnyTimes()

			mockMembershipRepo.EXPECT().FindByUserAndGroup(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, uid user.UserID, _ group.GroupID) (membership.Membership, error) {
					if uid == operatorID {
						return operatorMembership, tt.operatorFind
					}
					if tt.targetExists {
						return mocksmembership.NewMockMembership(ctrl), nil
					}
					return nil, membership.ErrMembershipNotFound
				}).AnyTimes()
			mockMembershipRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(tt.createErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.AddMember(context.Background(), groupID, targetID, tt.grantRole)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestRemoveMember(t *testing.T) {
	t.Parallel()
	operatorID := user.NewUserID()
	otherID := user.NewUserID()
	groupID := group.NewGroupID()

	tests := []struct {
		name         string
		success      bool
		targetID     user.UserID
		operatorRole membership.Role
		targetRole   membership.Role
		operatorFind error
		targetFind   error
		ownerCount   int
		deleteErr    error
	}{
		{"success self leave", true, operatorID, membership.RoleMember, membership.RoleMember, nil, nil, 0, nil},
		{"success admin removes member", true, otherID, membership.RoleAdmin, membership.RoleMember, nil, nil, 0, nil},
		{"failure permission denied", false, otherID, membership.RoleMember, membership.RoleMember, nil, nil, 0, nil},
		{"failure last owner", false, operatorID, membership.RoleOwner, membership.RoleOwner, nil, nil, 1, nil},
		{"success owner leaves when others exist", true, operatorID, membership.RoleOwner, membership.RoleOwner, nil, nil, 2, nil},
		{"failure operator not member", false, otherID, membership.RoleAdmin, membership.RoleMember, membership.ErrMembershipNotFound, nil, 0, nil},
		{"failure target not found", false, otherID, membership.RoleAdmin, membership.RoleMember, nil, membership.ErrMembershipNotFound, 0, nil},
		{"failure delete error", false, otherID, membership.RoleAdmin, membership.RoleMember, nil, nil, 0, errors.New("delete")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, operatorID, validIdentityID, nil, nil)

			operatorMembership := mocksmembership.NewMockMembership(ctrl)
			operatorMembership.EXPECT().Role().Return(tt.operatorRole).AnyTimes()
			targetMembership := mocksmembership.NewMockMembership(ctrl)
			targetMembership.EXPECT().Role().Return(tt.targetRole).AnyTimes()

			mockMembershipRepo.EXPECT().FindByUserAndGroup(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, uid user.UserID, _ group.GroupID) (membership.Membership, error) {
					if uid == operatorID && tt.targetID != operatorID {
						return operatorMembership, tt.operatorFind
					}
					return targetMembership, tt.targetFind
				}).AnyTimes()
			mockMembershipRepo.EXPECT().CountOwners(gomock.Any(), gomock.Any()).Return(tt.ownerCount, nil).AnyTimes()
			mockMembershipRepo.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.deleteErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			err := u.RemoveMember(context.Background(), groupID, tt.targetID)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestUpdateMemberRole(t *testing.T) {
	t.Parallel()
	operatorID := user.NewUserID()
	targetID := user.NewUserID()
	groupID := group.NewGroupID()

	tests := []struct {
		name         string
		success      bool
		operatorRole membership.Role
		targetRole   membership.Role
		newRole      membership.Role
		operatorFind error
		targetFind   error
		ownerCount   int
		updateErr    error
	}{
		{"success owner promotes member", true, membership.RoleOwner, membership.RoleMember, membership.RoleAdmin, nil, nil, 0, nil},
		{"failure not owner", false, membership.RoleAdmin, membership.RoleMember, membership.RoleAdmin, nil, nil, 0, nil},
		{"failure last owner downgrade", false, membership.RoleOwner, membership.RoleOwner, membership.RoleMember, nil, nil, 1, nil},
		{"success downgrade when others exist", true, membership.RoleOwner, membership.RoleOwner, membership.RoleMember, nil, nil, 2, nil},
		{"failure operator not member", false, membership.RoleOwner, membership.RoleMember, membership.RoleAdmin, membership.ErrMembershipNotFound, nil, 0, nil},
		{"failure target not found", false, membership.RoleOwner, membership.RoleMember, membership.RoleAdmin, nil, membership.ErrMembershipNotFound, 0, nil},
		{"failure update error", false, membership.RoleOwner, membership.RoleMember, membership.RoleAdmin, nil, nil, 0, errors.New("update")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, operatorID, validIdentityID, nil, nil)

			operatorMembership := mocksmembership.NewMockMembership(ctrl)
			operatorMembership.EXPECT().Role().Return(tt.operatorRole).AnyTimes()
			targetMembership := mocksmembership.NewMockMembership(ctrl)
			targetMembership.EXPECT().Role().Return(tt.targetRole).AnyTimes()
			targetMembership.EXPECT().ChangeRole(gomock.Any()).AnyTimes()

			mockMembershipRepo.EXPECT().FindByUserAndGroup(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
				func(_ context.Context, uid user.UserID, _ group.GroupID) (membership.Membership, error) {
					if uid == operatorID {
						return operatorMembership, tt.operatorFind
					}
					return targetMembership, tt.targetFind
				}).AnyTimes()
			mockMembershipRepo.EXPECT().CountOwners(gomock.Any(), gomock.Any()).Return(tt.ownerCount, nil).AnyTimes()
			mockMembershipRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(tt.updateErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.UpdateMemberRole(context.Background(), groupID, targetID, tt.newRole)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListMembers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		success   bool
		verifyErr error
		listErr   error
	}{
		{"success list members", true, nil, nil},
		{"failure verify token", false, errors.New("verify"), nil},
		{"failure list error", false, nil, errors.New("list")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, user.NewUserID(), validIdentityID, tt.verifyErr, nil)

			mockMembershipRepo.EXPECT().ListByGroupID(gomock.Any(), gomock.Any()).Return(nil, tt.listErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.ListMembers(context.Background(), group.NewGroupID())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestListMyGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		success      bool
		verifyErr    error
		listErr      error
		findByIDsErr error
	}{
		{"success list my groups", true, nil, nil, nil},
		{"failure verify token", false, errors.New("verify"), nil, nil},
		{"failure list error", false, nil, errors.New("list"), nil},
		{"failure find groups error", false, nil, nil, errors.New("find groups error")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuth := mocksappauth.NewMockAuthService(ctrl)
			mockUserRepo := mocksuser.NewMockUserRepository(ctrl)
			mockMembershipRepo := mocksmembership.NewMockMembershipRepository(ctrl)
			mockGroupRepo := mocksgroup.NewMockGroupRepository(ctrl)
			setupAuth(ctrl, mockAuth, mockUserRepo, user.NewUserID(), validIdentityID, tt.verifyErr, nil)

			gid := group.NewGroupID()
			mockMembership := mocksmembership.NewMockMembership(ctrl)
			mockMembership.EXPECT().GroupID().Return(gid).AnyTimes()
			mockMembership.EXPECT().Role().Return(membership.RoleOwner).AnyTimes()
			mockGroup := mocksgroup.NewMockGroup(ctrl)
			mockGroup.EXPECT().ID().Return(gid).AnyTimes()

			memberships := []membership.Membership{mockMembership}
			if tt.listErr != nil {
				memberships = nil
			}
			mockMembershipRepo.EXPECT().ListByUserID(gomock.Any(), gomock.Any()).Return(memberships, tt.listErr).AnyTimes()
			mockGroupRepo.EXPECT().FindByIDs(gomock.Any(), gomock.Any()).Return([]group.Group{mockGroup}, tt.findByIDsErr).AnyTimes()

			u := NewGroupUsecase(mockAuth, mockGroupRepo, mockMembershipRepo, mockUserRepo)

			_, err := u.ListMyGroups(context.Background())
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}
