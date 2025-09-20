package postauthlogin

import (
	"errors"
	"net/http"

	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/decorator"
	"github.com/Mafit1/notes-app/internal/api/common/refresh"
	"github.com/Mafit1/notes-app/internal/service/auth"
	"github.com/Mafit1/notes-app/internal/service/users"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type handler struct {
	authService auth.Service
}

func New(authService auth.Service) api.Handler {
	return decorator.NewBindAndValidate(&handler{authService})
}

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type Response struct {
	UserID      uuid.UUID `json:"id" validate:"required"`
	AccessToken string    `json:"accessToken" validate:"required"`
}

func (h *handler) Handle(c echo.Context, in Request) error {
	authData, err := h.authService.Login(
		c.Request().Context(),
		auth.LoginIn{
			Email:    in.Email,
			Password: in.Password,
		},
	)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if errors.Is(err, auth.ErrInvalidCredentials) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     refresh.RefreshTokenCookie,
		Value:    authData.RefreshData.RefreshToken,
		Expires:  authData.RefreshData.ExpiresAt,
		Path:     "/",
		HttpOnly: true,
	})

	return c.JSON(http.StatusAccepted, Response{
		UserID:      authData.UserData.ID,
		AccessToken: authData.AccessToken,
	})
}
