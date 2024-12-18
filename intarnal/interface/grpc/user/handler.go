package user

import (
	"context"
	"fmt"
	"user/intarnal/application/user"
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
	if err := h.userService.CreateUser(req.GetName(), req.GetEmail()); err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		Message: fmt.Sprintf("Create User, %s!", req.GetName()),
	}, nil
}
