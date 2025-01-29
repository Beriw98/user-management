package middleware

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

func NewLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			slog.With("method", c.Request().Method, "request", c.Request().RequestURI).Info("request received")

			return next(c)

		}
	}
}
