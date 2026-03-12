package handler

import (
	"net/http"
	"strconv"
	apimw "todo-api/middleware"
	"todo-api/service"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler(svc *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) GetAll(c echo.Context) error {
	userID := apimw.GetUserID(c)
	categories, err := h.svc.GetAll(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) Create(c echo.Context) error {
	userID := apimw.GetUserID(c)
	var body struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&body); err != nil || body.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	category, err := h.svc.Create(userID, body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	userID := apimw.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&body); err != nil || body.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	category, err := h.svc.Update(userID, id, body.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	userID := apimw.GetUserID(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(userID, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *CategoryHandler) Reorder(c echo.Context) error {
	userID := apimw.GetUserID(c)
	var body struct {
		IDs []int `json:"ids"`
	}
	if err := c.Bind(&body); err != nil || len(body.IDs) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := h.svc.Reorder(userID, body.IDs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
