package user

import (
	"context"
	"errors"
	"log"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userv1 "github.com/qkitzero/user-service/gen/go/user/v1"
	appuser "github.com/qkitzero/user-service/internal/application/user"
	domainuser "github.com/qkitzero/user-service/internal/domain/user"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	userUsecase appuser.UserUsecase
}

func NewUserHandler(
	userUsecase appuser.UserUsecase,
) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	displayName, err := domainuser.NewDisplayName(req.GetDisplayName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	birthDate, err := domainuser.NewBirthDate(req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := h.userUsecase.CreateUser(ctx, displayName, birthDate)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		log.Printf("CreateUser: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &userv1.CreateUserResponse{
		UserId: user.ID().String(),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	user, err := h.userUsecase.GetUser(ctx)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, domainuser.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("GetUser: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
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
	displayName, err := domainuser.NewDisplayName(req.GetDisplayName())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	birthDate, err := domainuser.NewBirthDate(req.GetBirthDate().GetYear(), req.GetBirthDate().GetMonth(), req.GetBirthDate().GetDay())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := h.userUsecase.UpdateUser(ctx, displayName, birthDate)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		if errors.Is(err, domainuser.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		log.Printf("UpdateUser: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &userv1.UpdateUserResponse{
		UserId:      user.ID().String(),
		DisplayName: user.DisplayName().String(),
		BirthDate: &date.Date{
			Year:  int32(user.BirthDate().Year()),
			Month: int32(user.BirthDate().Month()),
			Day:   int32(user.BirthDate().Day()),
		},
	}, nil
}
