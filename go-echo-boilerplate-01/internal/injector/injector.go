package injector

import (
	"github.com/joisandresky/go-echo-mysql-boilerplate/database"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/handler"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/middleware"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/repository"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/routes"
	"github.com/labstack/echo/v4"
)

type Injector interface {
	InjectModules()
}

type injector struct {
	conn   database.DatabaseProviderConnection
	server *echo.Echo
}

func NewInjector(conn database.DatabaseProviderConnection, server *echo.Echo) Injector {
	return &injector{conn, server}
}

func (i *injector) InjectModules() {
	authMw := middleware.NewAuthMiddleware()

	humanRepo := repository.NewHumanRepository(i.conn)
	humanHandler := handler.NewHumanHandler(humanRepo)
	humanRoutes := routes.NewHumanRoutes(humanHandler)

	humanRoutes.Install(i.server, authMw)
}
