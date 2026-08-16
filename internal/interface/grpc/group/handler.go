package group

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	groupv1 "github.com/qkitzero/user-service/gen/go/group/v1"
	appgroup "github.com/qkitzero/user-service/internal/application/group"
	domaingroup "github.com/qkitzero/user-service/internal/domain/group"
	"github.com/qkitzero/user-service/internal/domain/membership"
	"github.com/qkitzero/user-service/internal/domain/user"
)

type GroupHandler struct {
	groupv1.UnimplementedGroupServiceServer
	groupUsecase appgroup.GroupUsecase
}

func NewGroupHandler(
	groupUsecase appgroup.GroupUsecase,
) *GroupHandler {
	return &GroupHandler{
		groupUsecase: groupUsecase,
	}
}

func mapError(err error, op string) error {
	if _, ok := status.FromError(err); ok && status.Code(err) != codes.Unknown {
		return err
	}
	switch {
	case errors.Is(err, domaingroup.ErrGroupNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domaingroup.ErrCircularReference):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domaingroup.ErrAlreadyChild):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, user.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, membership.ErrMembershipNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, membership.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, membership.ErrAlreadyMember):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, membership.ErrLastOwnerCannotLeave):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		log.Printf("%s: internal error: %v", op, err)
		return status.Error(codes.Internal, "internal error")
	}
}

func toProtoGroup(g domaingroup.Group) *groupv1.Group {
	return &groupv1.Group{
		GroupId: g.ID().String(),
		Name:    g.Name().String(),
	}
}

func toProtoMember(m appgroup.Member) *groupv1.Member {
	return &groupv1.Member{
		UserId:      m.Membership.UserID().String(),
		Role:        m.Membership.Role().String(),
		DisplayName: m.DisplayName.String(),
	}
}

func (h *GroupHandler) CreateGroup(ctx context.Context, req *groupv1.CreateGroupRequest) (*groupv1.CreateGroupResponse, error) {
	name, err := domaingroup.NewGroupName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	group, err := h.groupUsecase.CreateGroup(ctx, name)
	if err != nil {
		return nil, mapError(err, "CreateGroup")
	}

	return &groupv1.CreateGroupResponse{
		GroupId: group.ID().String(),
	}, nil
}

func (h *GroupHandler) GetGroup(ctx context.Context, req *groupv1.GetGroupRequest) (*groupv1.GetGroupResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	group, err := h.groupUsecase.GetGroup(ctx, groupID)
	if err != nil {
		return nil, mapError(err, "GetGroup")
	}

	return &groupv1.GetGroupResponse{
		Group: toProtoGroup(group),
	}, nil
}

func (h *GroupHandler) UpdateGroup(ctx context.Context, req *groupv1.UpdateGroupRequest) (*groupv1.UpdateGroupResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	name, err := domaingroup.NewGroupName(req.GetName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	group, err := h.groupUsecase.UpdateGroup(ctx, groupID, name)
	if err != nil {
		return nil, mapError(err, "UpdateGroup")
	}

	return &groupv1.UpdateGroupResponse{
		Group: toProtoGroup(group),
	}, nil
}

func (h *GroupHandler) DeleteGroup(ctx context.Context, req *groupv1.DeleteGroupRequest) (*groupv1.DeleteGroupResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.groupUsecase.DeleteGroup(ctx, groupID); err != nil {
		return nil, mapError(err, "DeleteGroup")
	}

	return &groupv1.DeleteGroupResponse{}, nil
}

func (h *GroupHandler) ListChildGroups(ctx context.Context, req *groupv1.ListChildGroupsRequest) (*groupv1.ListChildGroupsResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	groups, err := h.groupUsecase.ListChildGroups(ctx, groupID)
	if err != nil {
		return nil, mapError(err, "ListChildGroups")
	}

	protoGroups := make([]*groupv1.Group, 0, len(groups))
	for _, g := range groups {
		protoGroups = append(protoGroups, toProtoGroup(g))
	}

	return &groupv1.ListChildGroupsResponse{
		Groups: protoGroups,
	}, nil
}

func (h *GroupHandler) ListParentGroups(ctx context.Context, req *groupv1.ListParentGroupsRequest) (*groupv1.ListParentGroupsResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	groups, err := h.groupUsecase.ListParentGroups(ctx, groupID)
	if err != nil {
		return nil, mapError(err, "ListParentGroups")
	}

	protoGroups := make([]*groupv1.Group, 0, len(groups))
	for _, g := range groups {
		protoGroups = append(protoGroups, toProtoGroup(g))
	}

	return &groupv1.ListParentGroupsResponse{
		Groups: protoGroups,
	}, nil
}

func (h *GroupHandler) AddChildGroup(ctx context.Context, req *groupv1.AddChildGroupRequest) (*groupv1.AddChildGroupResponse, error) {
	parentID, err := domaingroup.NewGroupIDFromString(req.GetParentGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	childID, err := domaingroup.NewGroupIDFromString(req.GetChildGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.groupUsecase.AddChildGroup(ctx, parentID, childID); err != nil {
		return nil, mapError(err, "AddChildGroup")
	}

	return &groupv1.AddChildGroupResponse{}, nil
}

func (h *GroupHandler) RemoveChildGroup(ctx context.Context, req *groupv1.RemoveChildGroupRequest) (*groupv1.RemoveChildGroupResponse, error) {
	parentID, err := domaingroup.NewGroupIDFromString(req.GetParentGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	childID, err := domaingroup.NewGroupIDFromString(req.GetChildGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.groupUsecase.RemoveChildGroup(ctx, parentID, childID); err != nil {
		return nil, mapError(err, "RemoveChildGroup")
	}

	return &groupv1.RemoveChildGroupResponse{}, nil
}

func (h *GroupHandler) AddMember(ctx context.Context, req *groupv1.AddMemberRequest) (*groupv1.AddMemberResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := user.NewUserIDFromString(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	role, err := membership.NewRole(req.GetRole())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := h.groupUsecase.AddMember(ctx, groupID, userID, role)
	if err != nil {
		return nil, mapError(err, "AddMember")
	}

	return &groupv1.AddMemberResponse{Member: toProtoMember(m)}, nil
}

func (h *GroupHandler) UpdateMemberRole(ctx context.Context, req *groupv1.UpdateMemberRoleRequest) (*groupv1.UpdateMemberRoleResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := user.NewUserIDFromString(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	role, err := membership.NewRole(req.GetRole())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := h.groupUsecase.UpdateMemberRole(ctx, groupID, userID, role)
	if err != nil {
		return nil, mapError(err, "UpdateMemberRole")
	}

	return &groupv1.UpdateMemberRoleResponse{Member: toProtoMember(m)}, nil
}

func (h *GroupHandler) RemoveMember(ctx context.Context, req *groupv1.RemoveMemberRequest) (*groupv1.RemoveMemberResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	userID, err := user.NewUserIDFromString(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.groupUsecase.RemoveMember(ctx, groupID, userID); err != nil {
		return nil, mapError(err, "RemoveMember")
	}

	return &groupv1.RemoveMemberResponse{}, nil
}

func (h *GroupHandler) ListMembers(ctx context.Context, req *groupv1.ListMembersRequest) (*groupv1.ListMembersResponse, error) {
	groupID, err := domaingroup.NewGroupIDFromString(req.GetGroupId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	members, err := h.groupUsecase.ListMembers(ctx, groupID)
	if err != nil {
		return nil, mapError(err, "ListMembers")
	}

	protoMembers := make([]*groupv1.Member, 0, len(members))
	for _, m := range members {
		protoMembers = append(protoMembers, toProtoMember(m))
	}

	return &groupv1.ListMembersResponse{Members: protoMembers}, nil
}

func (h *GroupHandler) ListMyGroups(ctx context.Context, _ *groupv1.ListMyGroupsRequest) (*groupv1.ListMyGroupsResponse, error) {
	groupMemberships, err := h.groupUsecase.ListMyGroups(ctx)
	if err != nil {
		return nil, mapError(err, "ListMyGroups")
	}

	groups := make([]*groupv1.GroupMembership, 0, len(groupMemberships))
	for _, gm := range groupMemberships {
		groups = append(groups, &groupv1.GroupMembership{
			Group: toProtoGroup(gm.Group),
			Role:  gm.Role.String(),
		})
	}

	return &groupv1.ListMyGroupsResponse{Groups: groups}, nil
}
