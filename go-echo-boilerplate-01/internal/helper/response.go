package helper

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Code       int         `json:"code"`
	Errors     interface{} `json:"errors,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

func CreatedResponse(c echo.Context, data Response) error {
	data.Success = true
	data.Code = http.StatusCreated

	return c.JSON(http.StatusCreated, data)
}

func OkResponse(c echo.Context, data Response) error {
	data.Success = true
	data.Code = http.StatusOK

	return c.JSON(http.StatusOK, data)
}

func NotFoundResponse(c echo.Context, data Response) error {
	data.Success = false
	data.Code = http.StatusNotFound

	return c.JSON(http.StatusNotFound, data)
}

func ForbiddenResponse(c echo.Context, data Response) error {
	data.Success = false
	data.Code = http.StatusForbidden

	return c.JSON(http.StatusForbidden, data)
}

func BadRequestResponse(c echo.Context, data Response) error {
	data.Success = false
	data.Code = http.StatusBadRequest

	return c.JSON(http.StatusBadRequest, data)
}

func UnprocResponse(c echo.Context, data Response) error {
	data.Success = false
	data.Code = http.StatusUnprocessableEntity

	return c.JSON(http.StatusUnprocessableEntity, data)
}

func ServerErrorResponse(c echo.Context, data Response) error {
	data.Success = false
	data.Code = http.StatusInternalServerError

	return c.JSON(http.StatusInternalServerError, data)
}

func UnauthorizedResponse(c echo.Context, data Response) error {
	data.Success = false
	data.Code = http.StatusUnauthorized

	return c.JSON(http.StatusUnauthorized, data)
}
