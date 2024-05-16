package handler

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/entity"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/helper"
	"github.com/joisandresky/go-echo-mysql-boilerplate/internal/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

type HumanHandler interface {
	GetAllHuman(c echo.Context) error
	GetHuman(c echo.Context) error
	CreateHuman(c echo.Context) error
	UpdateHuman(c echo.Context) error
	DeleteHuman(c echo.Context) error
}

type humanHandler struct {
	humanRepo repository.HumanRepository
}

func NewHumanHandler(humanRepo repository.HumanRepository) HumanHandler {
	return &humanHandler{humanRepo}
}

func (h *humanHandler) GetAllHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	humans, err := h.humanRepo.GetAll(ctx)
	if err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Failed to Get All Humant",
			Errors:  err.Error(),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Data: humans,
	})
}

func (h *humanHandler) GetHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	id, _ := strconv.Atoi(paramId)
	human, err := h.humanRepo.Show(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return helper.NotFoundResponse(c, helper.Response{
				Message: fmt.Sprintf("Human With [ID: %d] Not Found", id),
			})
		}

		return helper.BadRequestResponse(c, helper.Response{
			Message: fmt.Sprintf("Failed to Get Human With [ID: %d]", id),
			Errors:  err.Error(),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Data: human,
	})
}

func (h *humanHandler) CreateHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	var human entity.Human
	if err := c.Bind(&human); err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Failed to Get Value from Request Body",
			Errors:  err.Error(),
		})
	}

	if err := h.humanRepo.Store(ctx, human); err != nil {
		return helper.UnprocResponse(c, helper.Response{
			Message: "Failed to Save Human!",
			Errors:  err.Error(),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Message: "Human Saved!",
	})
}

func (h *humanHandler) UpdateHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	id, _ := strconv.Atoi(paramId)

	var human entity.Human
	if err := c.Bind(&human); err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: "Failed to Get Value from Request Body",
			Errors:  err.Error(),
		})
	}

	affected, err := h.humanRepo.Update(ctx, id, human)
	if err != nil {
		return helper.UnprocResponse(c, helper.Response{
			Message: "Failed to Save Human!",
			Errors:  err.Error(),
		})
	}

	if affected == 0 {
		return helper.NotFoundResponse(c, helper.Response{
			Message: fmt.Sprintf("Human With [ID: %d] Not Found", id),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Message: "Human Saved",
	})
}

func (h *humanHandler) DeleteHuman(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 10*time.Second)
	defer cancel()

	paramId := c.Param("id")
	id, _ := strconv.Atoi(paramId)

	affected, err := h.humanRepo.Delete(ctx, id)
	if err != nil {
		return helper.BadRequestResponse(c, helper.Response{
			Message: fmt.Sprintf("Failed to Delete Human With [ID: %d]", id),
			Errors:  err.Error(),
		})
	}

	if affected == 0 {
		return helper.NotFoundResponse(c, helper.Response{
			Message: fmt.Sprintf("Human With [ID: %d] Not Found", id),
		})
	}

	return helper.OkResponse(c, helper.Response{
		Message: "Human Deleted!",
	})
}
