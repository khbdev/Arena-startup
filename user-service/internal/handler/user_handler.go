package handler

import (
	"context"
	domain "user-service/internal/domen"
	"user-service/internal/usecase"

	userpb "github.com/khbdev/arena-startup-proto/proto/user"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	userpb.UnimplementedUserServiceServer
}

// NewUserHandler - constructor
func NewUserHandler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: u,
	}
}

// GetUserByTelegramId - proto rpc implementation
func (h *UserHandler) GetUserByTelegramId(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := h.userUsecase.GetUserByTelegramID(req.TelegramId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &userpb.GetUserResponse{User: nil}, nil
	}

	return &userpb.GetUserResponse{
		User: &userpb.User{
			TelegramId: user.TelegramID,
			Role:       user.Role,
			FirstName:  user.FirstName,
		},
	}, nil
}

// CreateUser - proto rpc implementation
func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	domainUser := &domain.User{
		TelegramID: req.TelegramId,
		Role:       req.Role,
		FirstName:  req.FirstName,
	}

	user, err := h.userUsecase.CreateUser(domainUser)
	if err != nil {
		return nil, err
	}

	return &userpb.CreateUserResponse{
		User: &userpb.User{
			TelegramId: user.TelegramID,
			Role:       user.Role,
			FirstName:  user.FirstName,
		},
	}, nil
}
