package user

import (
	"context"
	"user/internal/application/user"
	"user/pb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user, err := h.userService.CreateUser(req.GetDisplayName(), req.GetEmail())
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		UserId: user.ID().String(),
	}, nil
}
