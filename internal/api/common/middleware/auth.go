package middleware

import (
	"net/http"
	"strings"

	"github.com/Mafit1/notes-app/internal/models"
	"github.com/Mafit1/notes-app/internal/service/jwtservice"
	"github.com/labstack/echo/v4"
)

type AuthMW struct {
	jwtService jwtservice.Service
}

func NewAuthMW(jwtService jwtservice.Service) *AuthMW {
	return &AuthMW{jwtService}
}

func (m *AuthMW) Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header format")
			}

			claims, err := m.jwtService.ParseAccessToken(parts[1])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}

func (m *AuthMW) RequireRole(requiredRole models.RoleType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get("role").(models.RoleType)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "role not found in context")
			}

			if role != requiredRole {
				return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
			}

			return next(c)
		}
	}
}
