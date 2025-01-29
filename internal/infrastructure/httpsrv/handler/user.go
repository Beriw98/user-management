package handler

import (
	"github.com/Beriw98/user-management/ent"
	"github.com/Beriw98/user-management/internal/infrastructure/database/entity"
	"github.com/labstack/echo/v4"
)

type UserPostRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserPostResponse struct {
	ID string `json:"id"`
}

type UserPostHandler struct {
	userClient *ent.UserClient
}

func NewUserPostHandler(userRepository UserRepositoryInterface) *UserPostHandler {
	return &UserPostHandler{
		userRepository: userRepository,
	}
}

func (h *UserPostHandler) Handle(ec echo.Context) (*UserPostResponse, error) {
	ctx := ec.Request().Context()

	var req *UserPostRequest
	if err := ec.Bind(&req); err != nil {
		return nil, err
	}

	if err := ec.Validate(req); err != nil {
		return nil, err
	}

	h.userClient.Create().
		SetName(req.Name).
		SetSurname(req.Surname).
		SetEmail(req.Email).
		SetPassword(req.Password).Save(ctx)

	if _, err := h.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}
}
