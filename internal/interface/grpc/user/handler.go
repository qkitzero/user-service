package user

import (
	"context"

	"github.com/qkitzero/user/internal/application/auth"
	"github.com/qkitzero/user/internal/application/user"
	"github.com/qkitzero/user/pb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	authUsecase auth.AuthUsecase
	userUsecase user.UserUsecase
}

func NewUserHandler(
	authUsecase auth.AuthUsecase,
	userUsecase user.UserUsecase,
) *UserHandler {
	return &UserHandler{
		authUsecase: authUsecase,
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	userID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.CreateUser(userID, req.GetDisplayName(), req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		UserId: user.ID().String(),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	userID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.GetUser(userID)
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
	userID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.UpdateUser(userID, req.GetDisplayName())
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserResponse{
		UserId:      user.ID().String(),
		DisplayName: user.DisplayName().String(),
	}, nil
}
