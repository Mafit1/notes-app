package postauthlogout

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
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

func (h *handler) Handle(c echo.Context, in Request) error {
	refreshTokenCookie, err := c.Cookie("refreshToken")
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	err = h.authService.Logout(c.Request().Context(), refreshTokenCookie.Value)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}
		if errors.Is(err, auth.ErrCannotLogout) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	refreshTokenCookie.Expires = time.Unix(0, 0)
	refreshTokenCookie.HttpOnly = true
	refreshTokenCookie.Path = "/"
	c.SetCookie(refreshTokenCookie)

	return c.NoContent(http.StatusNoContent)
}
