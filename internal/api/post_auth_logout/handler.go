package postauthlogout

import (
	"errors"
	"net/http"
	"time"

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

func (h *handler) Handle(c echo.Context, in Request) error {
	refreshTokenCookie, err := refresh.ExtractRefreshTokenFromCookie(c)
	if err != nil {
		return c.NoContent(http.StatusNoContent)
	}

	err = h.authService.Logout(c.Request().Context(), refreshTokenCookie)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}
		if errors.Is(err, auth.ErrCannotLogout) {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	deleteCookie := &http.Cookie{
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	}
	c.SetCookie(deleteCookie)

	return c.NoContent(http.StatusNoContent)
}
