package refresh

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const RefreshTokenCookie = "refreshToken"

func ExtractRefreshTokenFromCookie(c echo.Context) (string, error) {
	cookie, err := c.Cookie(RefreshTokenCookie)
	if err != nil {
		return "", fmt.Errorf("refresh token cookie not found: %w", err)
	}
	return cookie.Value, nil
}
