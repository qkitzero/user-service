package user

import (
	"context"
	"errors"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "github.com/qkitzero/user/gen/go/user/v1"
	appauth "github.com/qkitzero/user/internal/application/auth"
	appuser "github.com/qkitzero/user/internal/application/user"
	domainuser "github.com/qkitzero/user/internal/domain/user"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	authUsecase appauth.AuthUsecase
	userUsecase appuser.UserUsecase
}

func NewUserHandler(
	authUsecase appauth.AuthUsecase,
	userUsecase appuser.UserUsecase,
) *UserHandler {
	return &UserHandler{
		authUsecase: authUsecase,
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	identityID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.CreateUser(identityID, req.GetDisplayName(), req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, err
	}

	return &userv1.CreateUserResponse{
		UserId: user.ID().String(),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	identityID, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.GetUser(identityID)
	if errors.Is(err, domainuser.ErrUserNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}
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
	_, err := h.authUsecase.VerifyToken(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userUsecase.UpdateUser(req.GetUserId(), req.GetDisplayName())
	if err != nil {
		return nil, err
	}

	return &userv1.UpdateUserResponse{
		UserId:      user.ID().String(),
		DisplayName: user.DisplayName().String(),
	}, nil
}
