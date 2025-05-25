package user

import (
	"context"

	"google.golang.org/genproto/googleapis/type/date"

	userv1 "github.com/qkitzero/user/gen/go/user/v1"
	"github.com/qkitzero/user/internal/application/auth"
	"github.com/qkitzero/user/internal/application/user"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
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

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	userID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.CreateUser(userID, req.GetDisplayName(), req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, err
	}

	return &userv1.CreateUserResponse{
		UserId: user.ID().String(),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	userID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.GetUser(userID)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserResponse{
		UserId:      user.ID().String(),
		DisplayName: user.DisplayName().String(),
		BirthDate: &date.Date{
			Year:  int32(user.BirthDate().Year()),
			Month: int32(user.BirthDate().Month()),
			Day:   int32(user.BirthDate().Day()),
		},
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	userID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.UpdateUser(userID, req.GetDisplayName())
	if err != nil {
		return nil, err
	}

	return &userv1.UpdateUserResponse{
		UserId:      user.ID().String(),
		DisplayName: user.DisplayName().String(),
	}, nil
}
