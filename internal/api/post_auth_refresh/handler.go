package postauthrefresh

import (
	"errors"
	"net/http"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
	"github.com/Mafit1/notes-app/internal/api/common/refresh"
	"github.com/Mafit1/notes-app/internal/service/auth"
	"github.com/labstack/echo/v4"
)

type handler struct {
	authService auth.Service
}

func New(authService auth.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{authService})
}

type Request struct{}

type Response struct {
	AccessToken string `json:"accessToken" validate:"required"`
}

func (h *handler) Handle(c echo.Context, _ Request) error {
	refreshCookie, err := c.Cookie("refreshToken")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
	}

	tokens, err := h.authService.RefreshTokens(c.Request().Context(), refreshCookie.Value)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     refresh.RefreshTokenCookie,
		Value:    tokens.RefreshToken,
		Expires:  tokens.RefreshExpiresAt,
		Path:     "/api/auth",
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, Response{
		AccessToken: tokens.AccessToken,
	})
}
