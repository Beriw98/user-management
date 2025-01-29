package httpsrv

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/Beriw98/user-management/internal/container"
	customvalidator "github.com/Beriw98/user-management/internal/infrastructure/httpsrv/handler/validator"
	"github.com/Beriw98/user-management/internal/infrastructure/httpsrv/middleware"
)

type requestValidator struct {
	Validator *validator.Validate
}

func (v *requestValidator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}

func NewRouter(ctr *container.Container) *echo.Echo {
	e := echo.New()

	v := &requestValidator{
		Validator: validator.New(),
	}

	_ = v.Validator.RegisterValidation(customvalidator.PasswordValidator, customvalidator.PasswordValidate)

	e.Validator = v
	e.HideBanner = true

	g := e.Group("/users", middleware.NewLoggerMiddleware())
	{
		g.POST("", ctr.UserHandler.Create)
		g.GET("", ctr.UserHandler.GetMany)
		g.GET("/:id", ctr.UserHandler.GetByID)
		g.PUT("/:id", ctr.UserHandler.Update)
		g.PATCH("/:id/password", ctr.UserHandler.UpdatePassword)
		g.DELETE("/:id", ctr.UserHandler.Delete)
	}

	return e
}
