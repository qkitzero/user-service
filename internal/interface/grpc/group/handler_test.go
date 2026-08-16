package group

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	groupv1 "github.com/qkitzero/user-service/gen/go/group/v1"
	appgroup "github.com/qkitzero/user-service/internal/application/group"
	"github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
	mocksappgroup "github.com/qkitzero/user-service/mocks/application/group"
	mocksgroup "github.com/qkitzero/user-service/mocks/domain/group"
	mocksmembership "github.com/qkitzero/user-service/mocks/domain/membership"
)

func validGroupID() string {
	return group.NewGroupID().String()
}

func validUserID() string {
	return user.NewUserID().String()
}

func TestCreateGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		groupName      string
		callUsecase    bool
		createGroupErr error
		wantCode       codes.Code
	}{
		{"success create group", "test group", true, nil, codes.OK},
		{"failure invalid name", "", false, nil, codes.InvalidArgument},
		{"failure usecase error", "test group", true, fmt.Errorf("create group error"), codes.Internal},
		{"failure status preserved", "test group", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockGroup := mocksgroup.NewMockGroup(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().CreateGroup(gomock.Any(), gomock.Any()).Return(mockGroup, tt.createGroupErr).Times(1)
				mockGroup.EXPECT().ID().Return(group.NewGroupID()).AnyTimes()
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.CreateGroupRequest{Name: tt.groupName}

			_, err := handler.CreateGroup(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestGetGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		groupID     string
		callUsecase bool
		getGroupErr error
		wantCode    codes.Code
	}{
		{"success get group", validGroupID(), true, nil, codes.OK},
		{"failure invalid group id", "invalid-uuid", false, nil, codes.InvalidArgument},
		{"failure group not found", validGroupID(), true, group.ErrGroupNotFound, codes.NotFound},
		{"failure usecase error", validGroupID(), true, fmt.Errorf("get group error"), codes.Internal},
		{"failure status preserved", validGroupID(), true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockGroup := mocksgroup.NewMockGroup(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().GetGroup(gomock.Any(), gomock.Any()).Return(mockGroup, tt.getGroupErr).Times(1)
				if tt.getGroupErr == nil {
					mockGroup.EXPECT().ID().Return(group.NewGroupID()).AnyTimes()
					mockGroup.EXPECT().Name().Return(group.GroupName("test group")).AnyTimes()
				}
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.GetGroupRequest{GroupId: tt.groupID}

			_, err := handler.GetGroup(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestUpdateGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		groupID        string
		groupName      string
		callUsecase    bool
		updateGroupErr error
		wantCode       codes.Code
	}{
		{"success update group", validGroupID(), "updated group", true, nil, codes.OK},
		{"failure invalid group id", "invalid-uuid", "updated group", false, nil, codes.InvalidArgument},
		{"failure invalid name", validGroupID(), "", false, nil, codes.InvalidArgument},
		{"failure group not found", validGroupID(), "updated group", true, group.ErrGroupNotFound, codes.NotFound},
		{"failure usecase error", validGroupID(), "updated group", true, fmt.Errorf("update group error"), codes.Internal},
		{"failure status preserved", validGroupID(), "updated group", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockGroup := mocksgroup.NewMockGroup(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().UpdateGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGroup, tt.updateGroupErr).Times(1)
				if tt.updateGroupErr == nil {
					mockGroup.EXPECT().ID().Return(group.NewGroupID()).AnyTimes()
					mockGroup.EXPECT().Name().Return(group.GroupName("updated group")).AnyTimes()
				}
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.UpdateGroupRequest{
				GroupId: tt.groupID,
				Name:    tt.groupName,
			}

			_, err := handler.UpdateGroup(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestDeleteGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		groupID        string
		callUsecase    bool
		deleteGroupErr error
		wantCode       codes.Code
	}{
		{"success delete group", validGroupID(), true, nil, codes.OK},
		{"failure invalid group id", "invalid-uuid", false, nil, codes.InvalidArgument},
		{"failure group not found", validGroupID(), true, group.ErrGroupNotFound, codes.NotFound},
		{"failure usecase error", validGroupID(), true, fmt.Errorf("delete group error"), codes.Internal},
		{"failure status preserved", validGroupID(), true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().DeleteGroup(gomock.Any(), gomock.Any()).Return(tt.deleteGroupErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.DeleteGroupRequest{GroupId: tt.groupID}

			_, err := handler.DeleteGroup(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestListChildGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		groupID      string
		callUsecase  bool
		listChildErr error
		wantCode     codes.Code
	}{
		{"success list child groups", validGroupID(), true, nil, codes.OK},
		{"failure invalid group id", "invalid-uuid", false, nil, codes.InvalidArgument},
		{"failure usecase error", validGroupID(), true, fmt.Errorf("list child groups error"), codes.Internal},
		{"failure status preserved", validGroupID(), true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			childGroupID := group.NewGroupID()
			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockGroup := mocksgroup.NewMockGroup(ctrl)
			if tt.callUsecase {
				var groups []group.Group
				if tt.listChildErr == nil {
					mockGroup.EXPECT().ID().Return(childGroupID).AnyTimes()
					mockGroup.EXPECT().Name().Return(group.GroupName("child group")).AnyTimes()
					groups = []group.Group{mockGroup}
				}
				mockUsecase.EXPECT().ListChildGroups(gomock.Any(), gomock.Any()).Return(groups, tt.listChildErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.ListChildGroupsRequest{GroupId: tt.groupID}

			resp, err := handler.ListChildGroups(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if len(resp.GetGroups()) != 1 {
					t.Fatalf("len = %d, want 1", len(resp.GetGroups()))
				}
				if got := resp.GetGroups()[0].GetGroupId(); got != childGroupID.String() {
					t.Errorf("group id = %v, want %v", got, childGroupID.String())
				}
				if got := resp.GetGroups()[0].GetName(); got != "child group" {
					t.Errorf("name = %v, want %v", got, "child group")
				}
			}
		})
	}
}

func TestListParentGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		groupID       string
		callUsecase   bool
		listParentErr error
		wantCode      codes.Code
	}{
		{"success list parent groups", validGroupID(), true, nil, codes.OK},
		{"failure invalid group id", "invalid-uuid", false, nil, codes.InvalidArgument},
		{"failure usecase error", validGroupID(), true, fmt.Errorf("list parent groups error"), codes.Internal},
		{"failure status preserved", validGroupID(), true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			parentGroupID := group.NewGroupID()
			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockGroup := mocksgroup.NewMockGroup(ctrl)
			if tt.callUsecase {
				var groups []group.Group
				if tt.listParentErr == nil {
					mockGroup.EXPECT().ID().Return(parentGroupID).AnyTimes()
					mockGroup.EXPECT().Name().Return(group.GroupName("parent group")).AnyTimes()
					groups = []group.Group{mockGroup}
				}
				mockUsecase.EXPECT().ListParentGroups(gomock.Any(), gomock.Any()).Return(groups, tt.listParentErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.ListParentGroupsRequest{GroupId: tt.groupID}

			resp, err := handler.ListParentGroups(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if len(resp.GetGroups()) != 1 {
					t.Fatalf("len = %d, want 1", len(resp.GetGroups()))
				}
				if got := resp.GetGroups()[0].GetGroupId(); got != parentGroupID.String() {
					t.Errorf("group id = %v, want %v", got, parentGroupID.String())
				}
				if got := resp.GetGroups()[0].GetName(); got != "parent group" {
					t.Errorf("name = %v, want %v", got, "parent group")
				}
			}
		})
	}
}

func TestAddChildGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		parentID    string
		childID     string
		callUsecase bool
		addErr      error
		wantCode    codes.Code
	}{
		{"success add child", validGroupID(), validGroupID(), true, nil, codes.OK},
		{"failure invalid parent id", "invalid", validGroupID(), false, nil, codes.InvalidArgument},
		{"failure invalid child id", validGroupID(), "invalid", false, nil, codes.InvalidArgument},
		{"failure circular reference", validGroupID(), validGroupID(), true, group.ErrCircularReference, codes.FailedPrecondition},
		{"failure already child", validGroupID(), validGroupID(), true, group.ErrAlreadyChild, codes.AlreadyExists},
		{"failure group not found", validGroupID(), validGroupID(), true, group.ErrGroupNotFound, codes.NotFound},
		{"failure internal error", validGroupID(), validGroupID(), true, fmt.Errorf("add child error"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().AddChildGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.addErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.AddChildGroupRequest{ParentGroupId: tt.parentID, ChildGroupId: tt.childID}

			_, err := handler.AddChildGroup(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestRemoveChildGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		parentID    string
		childID     string
		callUsecase bool
		removeErr   error
		wantCode    codes.Code
	}{
		{"success remove child", validGroupID(), validGroupID(), true, nil, codes.OK},
		{"failure invalid parent id", "invalid", validGroupID(), false, nil, codes.InvalidArgument},
		{"failure invalid child id", validGroupID(), "invalid", false, nil, codes.InvalidArgument},
		{"failure internal error", validGroupID(), validGroupID(), true, fmt.Errorf("remove child error"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().RemoveChildGroup(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.removeErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.RemoveChildGroupRequest{ParentGroupId: tt.parentID, ChildGroupId: tt.childID}

			_, err := handler.RemoveChildGroup(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestAddMember(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		groupID     string
		userID      string
		role        string
		callUsecase bool
		addErr      error
		wantCode    codes.Code
	}{
		{"success add member", validGroupID(), validUserID(), "member", true, nil, codes.OK},
		{"failure invalid group id", "invalid", validUserID(), "member", false, nil, codes.InvalidArgument},
		{"failure invalid user id", validGroupID(), "invalid", "member", false, nil, codes.InvalidArgument},
		{"failure invalid role", validGroupID(), validUserID(), "invalid", false, nil, codes.InvalidArgument},
		{"failure permission denied", validGroupID(), validUserID(), "member", true, membership.ErrPermissionDenied, codes.PermissionDenied},
		{"failure already member", validGroupID(), validUserID(), "member", true, membership.ErrAlreadyMember, codes.AlreadyExists},
		{"failure user not found", validGroupID(), validUserID(), "member", true, user.ErrUserNotFound, codes.NotFound},
		{"failure membership not found", validGroupID(), validUserID(), "member", true, membership.ErrMembershipNotFound, codes.NotFound},
		{"failure internal error", validGroupID(), validUserID(), "member", true, fmt.Errorf("add error"), codes.Internal},
		{"failure status preserved", validGroupID(), validUserID(), "member", true, status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockMembership := mocksmembership.NewMockMembership(ctrl)
			member := appgroup.Member{}
			if tt.callUsecase {
				if tt.addErr == nil {
					mockMembership.EXPECT().UserID().Return(user.NewUserID()).AnyTimes()
					mockMembership.EXPECT().Role().Return(membership.RoleMember).AnyTimes()
					member = appgroup.Member{Membership: mockMembership, DisplayName: user.DisplayName("test user")}
				}
				mockUsecase.EXPECT().AddMember(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(member, tt.addErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.AddMemberRequest{GroupId: tt.groupID, UserId: tt.userID, Role: tt.role}

			resp, err := handler.AddMember(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if got := resp.GetMember().GetDisplayName(); got != member.DisplayName.String() {
					t.Errorf("display name = %v, want %v", got, member.DisplayName.String())
				}
			}
		})
	}
}

func TestUpdateMemberRole(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		groupID     string
		userID      string
		role        string
		callUsecase bool
		updateErr   error
		wantCode    codes.Code
	}{
		{"success update role", validGroupID(), validUserID(), "admin", true, nil, codes.OK},
		{"failure invalid group id", "invalid", validUserID(), "admin", false, nil, codes.InvalidArgument},
		{"failure invalid user id", validGroupID(), "invalid", "admin", false, nil, codes.InvalidArgument},
		{"failure invalid role", validGroupID(), validUserID(), "invalid", false, nil, codes.InvalidArgument},
		{"failure permission denied", validGroupID(), validUserID(), "admin", true, membership.ErrPermissionDenied, codes.PermissionDenied},
		{"failure last owner", validGroupID(), validUserID(), "member", true, membership.ErrLastOwnerCannotLeave, codes.FailedPrecondition},
		{"failure internal error", validGroupID(), validUserID(), "admin", true, fmt.Errorf("update error"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockMembership := mocksmembership.NewMockMembership(ctrl)
			member := appgroup.Member{}
			if tt.callUsecase {
				if tt.updateErr == nil {
					mockMembership.EXPECT().UserID().Return(user.NewUserID()).AnyTimes()
					mockMembership.EXPECT().Role().Return(membership.RoleAdmin).AnyTimes()
					member = appgroup.Member{Membership: mockMembership, DisplayName: user.DisplayName("test user")}
				}
				mockUsecase.EXPECT().UpdateMemberRole(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(member, tt.updateErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.UpdateMemberRoleRequest{GroupId: tt.groupID, UserId: tt.userID, Role: tt.role}

			resp, err := handler.UpdateMemberRole(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if got := resp.GetMember().GetDisplayName(); got != member.DisplayName.String() {
					t.Errorf("display name = %v, want %v", got, member.DisplayName.String())
				}
			}
		})
	}
}

func TestRemoveMember(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		groupID     string
		userID      string
		callUsecase bool
		removeErr   error
		wantCode    codes.Code
	}{
		{"success remove member", validGroupID(), validUserID(), true, nil, codes.OK},
		{"failure invalid group id", "invalid", validUserID(), false, nil, codes.InvalidArgument},
		{"failure invalid user id", validGroupID(), "invalid", false, nil, codes.InvalidArgument},
		{"failure last owner", validGroupID(), validUserID(), true, membership.ErrLastOwnerCannotLeave, codes.FailedPrecondition},
		{"failure permission denied", validGroupID(), validUserID(), true, membership.ErrPermissionDenied, codes.PermissionDenied},
		{"failure internal error", validGroupID(), validUserID(), true, fmt.Errorf("remove error"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			if tt.callUsecase {
				mockUsecase.EXPECT().RemoveMember(gomock.Any(), gomock.Any(), gomock.Any()).Return(tt.removeErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.RemoveMemberRequest{GroupId: tt.groupID, UserId: tt.userID}

			_, err := handler.RemoveMember(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}

func TestListMembers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		groupID     string
		callUsecase bool
		listErr     error
		wantCode    codes.Code
	}{
		{"success list members", validGroupID(), true, nil, codes.OK},
		{"failure invalid group id", "invalid", false, nil, codes.InvalidArgument},
		{"failure permission denied", validGroupID(), true, membership.ErrPermissionDenied, codes.PermissionDenied},
		{"failure internal error", validGroupID(), true, fmt.Errorf("list error"), codes.Internal},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			memberUserID := user.NewUserID()
			memberDisplayName := user.DisplayName("test user")
			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockMembership := mocksmembership.NewMockMembership(ctrl)
			if tt.callUsecase {
				var members []appgroup.Member
				if tt.listErr == nil {
					mockMembership.EXPECT().UserID().Return(memberUserID).AnyTimes()
					mockMembership.EXPECT().Role().Return(membership.RoleMember).AnyTimes()
					members = []appgroup.Member{{Membership: mockMembership, DisplayName: memberDisplayName}}
				}
				mockUsecase.EXPECT().ListMembers(gomock.Any(), gomock.Any()).Return(members, tt.listErr).Times(1)
			}

			handler := NewGroupHandler(mockUsecase)

			req := &groupv1.ListMembersRequest{GroupId: tt.groupID}

			resp, err := handler.ListMembers(context.Background(), req)
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
			if tt.wantCode == codes.OK {
				if len(resp.GetMembers()) != 1 {
					t.Fatalf("len = %d, want 1", len(resp.GetMembers()))
				}
				if got := resp.GetMembers()[0].GetUserId(); got != memberUserID.String() {
					t.Errorf("user id = %v, want %v", got, memberUserID.String())
				}
				if got := resp.GetMembers()[0].GetRole(); got != membership.RoleMember.String() {
					t.Errorf("role = %v, want %v", got, membership.RoleMember.String())
				}
				if got := resp.GetMembers()[0].GetDisplayName(); got != memberDisplayName.String() {
					t.Errorf("display name = %v, want %v", got, memberDisplayName.String())
				}
			}
		})
	}
}

func TestListMyGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		listErr  error
		wantCode codes.Code
	}{
		{"success list my groups", nil, codes.OK},
		{"failure internal error", fmt.Errorf("list error"), codes.Internal},
		{"failure status preserved", status.Error(codes.Unauthenticated, "auth"), codes.Unauthenticated},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUsecase := mocksappgroup.NewMockGroupUsecase(ctrl)
			mockGroup := mocksgroup.NewMockGroup(ctrl)

			var result []appgroup.GroupMembership
			if tt.listErr == nil {
				mockGroup.EXPECT().ID().Return(group.NewGroupID()).AnyTimes()
				mockGroup.EXPECT().Name().Return(group.GroupName("test group")).AnyTimes()
				result = []appgroup.GroupMembership{{Group: mockGroup, Role: membership.RoleOwner}}
			}
			mockUsecase.EXPECT().ListMyGroups(gomock.Any()).Return(result, tt.listErr).Times(1)

			handler := NewGroupHandler(mockUsecase)

			_, err := handler.ListMyGroups(context.Background(), &groupv1.ListMyGroupsRequest{})
			if got := status.Code(err); got != tt.wantCode {
				t.Errorf("expected code %v, got %v (err=%v)", tt.wantCode, got, err)
			}
		})
	}
}
