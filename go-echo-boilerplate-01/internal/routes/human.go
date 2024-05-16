package routes

import (
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/handler"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/middleware"
	"github.com/labstack/echo/v4"
)

type humanRoutes struct {
	humanHandler handler.HumanHandler
}

func NewHumanRoutes(humanHandler handler.HumanHandler) Routes {
	return &humanRoutes{humanHandler}
}

func (r *humanRoutes) Install(server *echo.Echo, authMw middleware.AuthMiddleware) {
	v1 := server.Group("/api/v1/humans")
	v1.GET("", r.humanHandler.GetAllHuman)
	v1.GET("/:id", r.humanHandler.GetHuman)
	v1.POST("", r.humanHandler.CreateHuman)
	v1.PUT("/:id", r.humanHandler.UpdateHuman)
	v1.DELETE("/:id", r.humanHandler.DeleteHuman)

	/*
		if you want to use Auth Middleware("/internal/middleware/auth.go")
		just put middleware function in last params of echo http method
		eg: v1.GET("", r.humanHandler.GetAllHuman, authMw.IsAuthenticated)
	*/

	private := server.Group("/api/v1/private/humans")
	private.GET("", r.humanHandler.GetAllHuman, authMw.IsAuthenticated)
}
