package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Beriw98/user-management/internal/app/domain"
	"github.com/Beriw98/user-management/internal/infrastructure/httpsrv/handler/request"
	"github.com/Beriw98/user-management/internal/infrastructure/httpsrv/handler/response"
	customvalidator "github.com/Beriw98/user-management/internal/infrastructure/httpsrv/handler/validator"
)

type userRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
	GetMany(ctx context.Context, limit, offset int) ([]domain.User, error)
}

type UserHTTPHandler struct {
	userRepository userRepository
}

const (
	defaultLimit = "10"
	defaultPage  = "0"
)

func NewUserHTTPHandler(repository userRepository) *UserHTTPHandler {
	return &UserHTTPHandler{
		userRepository: repository,
	}
}

func (h *UserHTTPHandler) Create(ec echo.Context) error {
	ctx := ec.Request().Context()

	l := slog.Default().With("handler", "Create")

	var req *request.UserCreateRequest
	if err := ec.Bind(&req); err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ec.Validate(req); err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	if user != nil {
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	}

	id := xid.New().String()
	err = h.userRepository.Create(ctx, domain.User{
		ID:       id,
		Name:     req.Name,
		Surname:  req.Surname,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	return ec.JSON(http.StatusCreated, &response.UserIDResponse{ID: id})
}

func (h *UserHTTPHandler) GetByID(ec echo.Context) error {
	ctx := ec.Request().Context()

	l := slog.Default().With("handler", "GetByID")
	id := ec.Param("id")
	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	return ec.JSON(http.StatusOK, &response.UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Surname: user.Surname,
		Email:   user.Email,
	})
}

func (h *UserHTTPHandler) Update(ec echo.Context) error {
	ctx := ec.Request().Context()
	l := slog.Default().With("handler", "Update")
	id := ec.Param("id")

	var req *request.UserUpdateRequest
	if err := ec.Bind(&req); err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ec.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	user.Name = req.Name
	user.Surname = req.Surname
	user.Email = req.Email

	if err = h.userRepository.Update(ctx, *user); err != nil {
		return echo.ErrInternalServerError
	}

	return ec.JSON(http.StatusOK, &response.UserIDResponse{ID: id})
}

func (h *UserHTTPHandler) UpdatePassword(ec echo.Context) error {
	ctx := ec.Request().Context()
	l := slog.Default().With("handler", "UpdatePassword")

	id := ec.Param("id")

	var req *request.UserUpdatePasswordRequest
	if err := ec.Bind(&req); err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ec.Validate(req); err != nil {
		var vErr validator.ValidationErrors
		if errors.As(err, &vErr) {
			for _, e := range vErr {
				if e.Tag() == customvalidator.PasswordValidator {
					return echo.NewHTTPError(http.StatusBadRequest, customvalidator.ErrPasswordValidation)
				}
			}
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	user.Password = string(bytes)

	if err = h.userRepository.Update(ctx, *user); err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	return ec.JSON(http.StatusOK, &response.UserIDResponse{ID: id})
}

func (h *UserHTTPHandler) Delete(ec echo.Context) error {
	ctx := ec.Request().Context()
	l := slog.Default().With("handler", "Delete")

	id := ec.Param("id")
	user, err := h.userRepository.GetByID(ctx, id)
	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if err = h.userRepository.Delete(ctx, id); err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	return ec.NoContent(http.StatusNoContent)
}

func (h *UserHTTPHandler) GetMany(ec echo.Context) error {
	ctx := ec.Request().Context()
	l := slog.Default().With("handler", "GetMany")

	limit, page := ec.QueryParam("limit"), ec.QueryParam("page")
	if limit == "" {
		limit = defaultLimit
	}

	if page == "" {
		page = defaultPage
	}

	li, err := strconv.Atoi(limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "limit must be an integer")
	}

	p, err := strconv.Atoi(page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "offset must be an integer")
	}

	users, err := h.userRepository.GetMany(ctx, li, p)
	if err != nil {
		l.ErrorContext(ctx, err.Error())
		return echo.ErrInternalServerError
	}

	responseUsers := make([]response.UserResponse, 0)
	for _, user := range users {
		responseUsers = append(responseUsers, response.UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Surname: user.Surname,
			Email:   user.Email,
		})
	}

	return ec.JSON(http.StatusOK, responseUsers)
}
