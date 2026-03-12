package handler

import (
	"errors"
	"net/http"
	"time"
	"todo-api/service"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	svc *service.UserService
}

func NewAuthHandler(svc *service.UserService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil || body.Email == "" || body.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := h.svc.Register(body.Email, body.Password); err != nil {
		var ve service.ValidationError
		if errors.As(err, &ve) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusAccepted, map[string]string{"message": "確認メールを送信しました"})
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing token")
	}
	if err := h.svc.VerifyEmail(token); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "メール認証が完了しました"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil || body.Email == "" || body.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	accessToken, refreshToken, err := h.svc.Login(body.Email, body.Password)
	if err != nil {
		if err.Error() == "email not verified" {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   900,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth/refresh",
		MaxAge:   604800,
	})
	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		clearCookies(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
	}
	newAccess, newRefresh, err := h.svc.Refresh(cookie.Value)
	if err != nil {
		clearCookies(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    newAccess,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   900,
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    newRefresh,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth/refresh",
		MaxAge:   604800,
	})
	return c.NoContent(http.StatusOK)
}

func (h *AuthHandler) RequestPasswordReset(c echo.Context) error {
	var body struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&body); err != nil || body.Email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	_ = h.svc.RequestPasswordReset(body.Email)
	return c.JSON(http.StatusAccepted, map[string]string{"message": "パスワード再設定メールを送信しました"})
}

func (h *AuthHandler) ConfirmPasswordReset(c echo.Context) error {
	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.Bind(&body); err != nil || body.Token == "" || body.NewPassword == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	if err := h.svc.ConfirmPasswordReset(body.Token, body.NewPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "パスワードを変更しました"})
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err == nil && cookie.Value != "" {
		_ = h.svc.Logout(cookie.Value)
	}
	clearCookies(c)
	return c.NoContent(http.StatusOK)
}

func clearCookies(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/api/auth/refresh",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}
