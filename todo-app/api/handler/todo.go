package handler

import (
	"net/http"
	"strconv"
	"todo-api/service"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	svc *service.TodoService
}

func NewTodoHandler(svc *service.TodoService) *TodoHandler {
	return &TodoHandler{svc: svc}
}

func (h *TodoHandler) GetAll(c echo.Context) error {
	todos, err := h.svc.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Create(c echo.Context) error {
	var body struct {
		Text string `json:"text"`
	}
	if err := c.Bind(&body); err != nil || body.Text == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	todo, err := h.svc.Create(body.Text)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	var body struct {
		Done bool `json:"done"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	todo, err := h.svc.Update(id, body.Done)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
