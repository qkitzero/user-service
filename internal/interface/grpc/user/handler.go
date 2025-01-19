package user

import (
	"context"
	"user/internal/application/auth"
	"user/internal/application/user"
	"user/pb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	authService auth.AuthService
	userService user.UserService
}

func NewUserHandler(
	authService auth.AuthService,
	userService user.UserService,
) *UserHandler {
	return &UserHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID, err := h.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.CreateUser(userID, req.GetDisplayName())
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		UserId: user.ID().String(),
	}, nil
}
