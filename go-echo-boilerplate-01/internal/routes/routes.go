package routes

import (
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/middleware"
	"github.com/labstack/echo/v4"
)

type Routes interface {
	Install(server *echo.Echo, authMw middleware.AuthMiddleware)
}
