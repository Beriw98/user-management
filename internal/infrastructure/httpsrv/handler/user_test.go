package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Beriw98/user-management/internal/app/domain"
	"github.com/Beriw98/user-management/internal/infrastructure/httpsrv/handler"
	customvalidator "github.com/Beriw98/user-management/internal/infrastructure/httpsrv/handler/validator"
)

type requestValidator struct {
	Validator *validator.Validate
}

func (v *requestValidator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Create(ctx context.Context, user domain.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *repositoryMock) GetByID(ctx context.Context, id string) (*domain.User, error) {
	args := r.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (r *repositoryMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := r.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (r *repositoryMock) Update(ctx context.Context, user domain.User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *repositoryMock) Delete(ctx context.Context, id string) error {
	args := r.Called(ctx, id)
	return args.Error(0)
}

func (r *repositoryMock) GetMany(ctx context.Context, limit, offset int) ([]domain.User, error) {
	args := r.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func TestNewUserHTTPHandler(t *testing.T) {
	t.Run("NewUserHTTPHandler", func(t *testing.T) {
		h := handler.NewUserHTTPHandler(nil)
		assert.NotNil(t, h)
	})
}

func TestUserHTTPHandler_Create(t *testing.T) {
	rm := new(repositoryMock)
	h := handler.NewUserHTTPHandler(rm)
	e := echo.New()
	e.Validator = &requestValidator{
		Validator: validator.New(),
	}

	t.Run("Create", func(t *testing.T) {

		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		ec := e.NewContext(req, res)
		ctx := ec.Request().Context()
		user := domain.User{
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "password",
		}

		rm.On("GetByEmail", ctx, user.Email).Return(nil, nil).Once()
		rm.On("Create", ctx, mock.Anything).Run(func(args mock.Arguments) {
			arg := args.Get(1).(domain.User)
			assert.Equal(t, user.Name, arg.Name)
			assert.Equal(t, user.Surname, arg.Surname)
			assert.Equal(t, user.Email, arg.Email)
			assert.Equal(t, user.Password, arg.Password)
		}).Return(nil).Once()

		err := h.Create(ec)

		assert.NoError(t, err)

		rm.AssertExpectations(t)
	})

	t.Run("Bind error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		ec := e.NewContext(req, res)
		err := h.Create(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":null,"password":"password"}`

		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		ec := e.NewContext(req, res)

		err := h.Create(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("GetByEmail error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		ec := e.NewContext(req, res)
		ctx := ec.Request().Context()
		user := domain.User{
			Name:    "Test",
			Surname: "Test",
			Email:   "test@test.pl",
		}

		rm.On("GetByEmail", ctx, user.Email).Return(nil, assert.AnError).Once()

		err := h.Create(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Create error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ctx := ec.Request().Context()
		user := domain.User{
			Name:    "Test",
			Surname: "Test",
			Email:   "test@test.pl",
		}

		rm.On("GetByEmail", ctx, user.Email).Return(nil, nil).Once()
		rm.On("Create", ctx, mock.Anything).Return(assert.AnError).Once()

		err := h.Create(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Conflict", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`
		req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		ec := e.NewContext(req, res)
		ctx := ec.Request().Context()

		user := domain.User{
			Name:    "Test",
			Surname: "Test",
			Email:   "test@test.pl",
		}

		rm.On("GetByEmail", ctx, user.Email).Return(&user, nil).Once()

		err := h.Create(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusConflict, he.Code)

		rm.AssertExpectations(t)
	})
}

func TestUserHTTPHandler_Delete(t *testing.T) {
	rm := new(repositoryMock)
	h := handler.NewUserHTTPHandler(rm)
	e := echo.New()
	e.Validator = &requestValidator{
		Validator: validator.New(),
	}

	t.Run("Delete", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")
		ctx := ec.Request().Context()

		rm.On("GetByID", ctx, "1").Return(&domain.User{ID: "1"}, nil).Once()
		rm.On("Delete", ctx, "1").Return(nil).Once()

		err := h.Delete(ec)

		assert.NoError(t, err)

		rm.AssertExpectations(t)
	})

	t.Run("Delete error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		ctx := ec.Request().Context()

		rm.On("GetByID", ctx, "1").Return(&domain.User{ID: "1"}, nil).Once()
		rm.On("Delete", ctx, "1").Return(assert.AnError).Once()

		err := h.Delete(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		ctx := ec.Request().Context()
		rm.On("GetByID", ctx, "1").Return(nil, nil).Once()

		err := h.Delete(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusNotFound, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("User not found error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")
		ctx := ec.Request().Context()

		rm.On("GetByID", ctx, "1").Return(nil, assert.AnError).Once()

		err := h.Delete(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})
}

func TestUserHTTPHandler_GetByID(t *testing.T) {
	rm := new(repositoryMock)
	h := handler.NewUserHTTPHandler(rm)
	e := echo.New()
	e.Validator = &requestValidator{
		Validator: validator.New(),
	}

	t.Run("GetByID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")
		ctx := ec.Request().Context()

		user := domain.User{
			ID:      "1",
			Name:    "Test",
			Surname: "Test",
			Email:   "test@test.pl",
		}

		rm.On("GetByID", ctx, "1").Return(&user, nil).Once()

		err := h.GetByID(ec)

		assert.NoError(t, err)

		rm.AssertExpectations(t)
	})

	t.Run("GetByID error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")
		ctx := ec.Request().Context()

		rm.On("GetByID", ctx, "1").Return(nil, assert.AnError).Once()

		err := h.GetByID(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users/1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")
		ctx := ec.Request().Context()

		rm.On("GetByID", ctx, "1").Return(nil, nil).Once()

		err := h.GetByID(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusNotFound, he.Code)

		rm.AssertExpectations(t)
	})

}

func TestUserHTTPHandler_GetMany(t *testing.T) {
	rm := new(repositoryMock)
	h := handler.NewUserHTTPHandler(rm)
	e := echo.New()
	e.Validator = &requestValidator{
		Validator: validator.New(),
	}

	t.Run("GetMany", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ctx := ec.Request().Context()

		users := []domain.User{
			{
				ID:      "1",
				Name:    "Test",
				Surname: "Test",
				Email:   "test@test.pl",
			},
		}

		rm.On("GetMany", ctx, 10, 0).Return(users, nil).Once()

		err := h.GetMany(ec)

		assert.NoError(t, err)

		rm.AssertExpectations(t)
	})

	t.Run("GetMany error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ctx := ec.Request().Context()

		rm.On("GetMany", ctx, 10, 0).Return(nil, assert.AnError).Once()

		err := h.GetMany(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)

	})

	t.Run("GetMany with limit and page query", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users?limit=5&page=1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ctx := ec.Request().Context()

		users := []domain.User{
			{
				ID:      "1",
				Name:    "Test",
				Surname: "Test",
				Email:   "test@test.pl",
			},
		}

		rm.On("GetMany", ctx, 5, 1).Return(users, nil).Once()

		err := h.GetMany(ec)

		assert.NoError(t, err)

		rm.AssertExpectations(t)
	})

	t.Run("GetMany with invalid limit", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users?limit=invalid&page=1", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)

		err := h.GetMany(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("GetMany with invalid page", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/users?limit=5&page=invalid", nil)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)

		err := h.GetMany(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})
}

func TestUserHTTPHandler_Update(t *testing.T) {
	rm := new(repositoryMock)
	h := handler.NewUserHTTPHandler(rm)
	e := echo.New()
	e.Validator = &requestValidator{
		Validator: validator.New(),
	}

	t.Run("Update", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`

		req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		ctx := ec.Request().Context()
		user := domain.User{
			ID:       "1",
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "password",
		}

		rm.On("GetByID", ctx, "1").Return(&user, nil).Once()
		rm.On("Update", ctx, mock.Anything).Run(func(args mock.Arguments) {
			arg := args.Get(1).(domain.User)
			assert.Equal(t, user.Name, arg.Name)
			assert.Equal(t, user.Surname, arg.Surname)
			assert.Equal(t, user.Email, arg.Email)
			assert.Equal(t, user.Password, arg.Password)
		}).Return(nil).Once()

		err := h.Update(ec)

		assert.NoError(t, err)

		rm.AssertExpectations(t)
	})

	t.Run("Update not found", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`

		req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()

		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		ctx := ec.Request().Context()

		rm.On("GetByID", ctx, "1").Return(nil, nil).Once()

		err := h.Update(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusNotFound, he.Code)
	})

	t.Run("Update error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`

		req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		ctx := ec.Request().Context()
		user := domain.User{
			ID:       "1",
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "password",
		}

		rm.On("GetByID", ctx, "1").Return(&user, nil).Once()
		rm.On("Update", ctx, mock.Anything).Return(assert.AnError).Once()

		err := h.Update(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Bind error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"`
		req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		err := h.Update(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":null,"password":"password"}`

		req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		err := h.Update(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("GetByID error", func(t *testing.T) {
		body := `{"name":"Test","surname":"Test","email":"test@test.pl","password":"password"}`

		req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		ctx := ec.Request().Context()

		rm.On("GetByID", ctx, "1").Return(nil, assert.AnError).Once()

		err := h.Update(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})
}

func TestUserHTTPHandler_UpdatePassword(t *testing.T) {
	rm := new(repositoryMock)
	h := handler.NewUserHTTPHandler(rm)
	e := echo.New()
	v := &requestValidator{
		Validator: validator.New(),
	}

	_ = v.Validator.RegisterValidation("password", customvalidator.PasswordValidate)
	e.Validator = v

	t.Run("UpdatePassword", func(t *testing.T) {
		body := `{"password":"1Password."}`
		req, _ := http.NewRequest(http.MethodPut, "/users/1/password", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")
		ctx := ec.Request().Context()

		user := domain.User{
			ID:       "1",
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "1Password.",
		}

		rm.On("GetByID", ctx, "1").Return(&user, nil).Once()
		rm.On("Update", ctx, mock.Anything).Run(func(args mock.Arguments) {
			arg := args.Get(1).(domain.User)
			assert.Equal(t, user.Name, arg.Name)
			assert.Equal(t, user.Surname, arg.Surname)
			assert.Equal(t, user.Email, arg.Email)
			assert.Equal(t, user.Password, arg.Password)
		}).Return(nil).Once()

		err := h.UpdatePassword(ec)

		assert.NoError(t, err)

		rm.AssertExpectations(t)
	})

	t.Run("UpdatePassword error", func(t *testing.T) {
		body := `{"password":"1Password."}`
		req, _ := http.NewRequest(http.MethodPut, "/users/1/password", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")
		ctx := ec.Request().Context()

		user := domain.User{
			ID:       "1",
			Name:     "Test",
			Surname:  "Test",
			Email:    "test@test.pl",
			Password: "1Password.",
		}

		rm.On("GetByID", ctx, "1").Return(&user, nil).Once()
		rm.On("Update", ctx, mock.Anything).Return(assert.AnError).Once()

		err := h.UpdatePassword(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusInternalServerError, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Bind error", func(t *testing.T) {
		body := `{"password":"`
		req, _ := http.NewRequest(http.MethodPut, "/users/1/password", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		err := h.UpdatePassword(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Validation error", func(t *testing.T) {
		body := `{"password":null}`
		req, _ := http.NewRequest(http.MethodPut, "/users/1/password", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		err := h.UpdatePassword(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)

		rm.AssertExpectations(t)
	})

	t.Run("Password validation error", func(t *testing.T) {
		body := `{"password":"password"}`
		req, _ := http.NewRequest(http.MethodPut, "/users/1/password", bytes.NewReader([]byte(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		res := httptest.NewRecorder()
		ec := e.NewContext(req, res)
		ec.SetParamNames("id")
		ec.SetParamValues("1")

		err := h.UpdatePassword(ec)

		var he *echo.HTTPError
		assert.ErrorAs(t, err, &he)

		assert.Equal(t, http.StatusBadRequest, he.Code)
	})
}
