package middleware

import (
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/helper"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/token"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

type AuthMiddleware interface {
	IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc
}

type authMiddleware struct {
	// Some Dependecies here ... (eg: permission repo or role repo)
}

func NewAuthMiddleware() AuthMiddleware {
	return &authMiddleware{}
}

func (mw *authMiddleware) IsAuthenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		maker, err := token.NewPasetoMaker(viper.GetString("app.asymmetricKey"))
		if err != nil {
			return helper.UnauthorizedResponse(c, helper.Response{
				Message: "Unauthorized",
			})
		}

		if c.Request().Header.Get("Authorization") == "" {
			return helper.UnauthorizedResponse(c, helper.Response{
				Message: "Unauthorized",
			})
		}

		access_token := c.Request().Header.Get("Authorization")

		payload, err := maker.VerifyToken(access_token)

		if err != nil {
			return helper.UnauthorizedResponse(c, helper.Response{
				Message: "Unauthorized [Token Invalid]",
			})
		}

		c.Set("claim", payload)

		return next(c)
	}
}
