package user

import (
	"context"

	"github.com/qkitzero/user/internal/application/auth"
	"github.com/qkitzero/user/internal/application/user"
	"github.com/qkitzero/user/pb"
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

	user, err := h.userService.CreateUser(userID, req.GetDisplayName(), req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		UserId: user.ID().String(),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID, err := h.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		UserId:      user.ID().String(),
		DisplayName: user.DisplayName().String(),
		BirthDate: &pb.Date{
			Year:  int32(user.BirthDate().Year()),
			Month: int32(user.BirthDate().Month()),
			Day:   int32(user.BirthDate().Day()),
		},
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	userID, err := h.authService.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.UpdateUser(userID, req.GetDisplayName())
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		UserId:      user.ID().String(),
		DisplayName: user.DisplayName().String(),
	}, nil
}
