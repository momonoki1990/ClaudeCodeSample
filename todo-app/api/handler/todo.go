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
	var categoryID *int
	if s := c.QueryParam("category_id"); s != "" {
		id, err := strconv.Atoi(s)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid category_id")
		}
		categoryID = &id
	}
	todos, err := h.svc.GetAll(categoryID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) Create(c echo.Context) error {
	var body struct {
		Text       string `json:"text"`
		CategoryID *int   `json:"category_id"`
	}
	if err := c.Bind(&body); err != nil || body.Text == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	todo, err := h.svc.Create(body.Text, body.CategoryID)
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
		Text       string `json:"text"`
		Done       bool   `json:"done"`
		CategoryID *int   `json:"category_id"`
	}
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	todo, err := h.svc.Update(id, body.Text, body.Done, body.CategoryID)
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

func (h *TodoHandler) Reorder(c echo.Context) error {
	var body struct {
		IDs []int `json:"ids"`
	}
	if err := c.Bind(&body); err != nil || len(body.IDs) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := h.svc.Reorder(body.IDs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *TodoHandler) DeleteDone(c echo.Context) error {
	var categoryID *int
	if s := c.QueryParam("category_id"); s != "" {
		id, err := strconv.Atoi(s)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid category_id")
		}
		categoryID = &id
	}
	if err := h.svc.DeleteDone(categoryID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
