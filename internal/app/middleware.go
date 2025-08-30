package app

import "github.com/Mafit1/notes-app/internal/api/common/middleware"

func (app *App) AuthMW() *middleware.AuthMW {
	if app.authMW == nil {
		app.authMW = middleware.New(app.JwtService())
	}
	return app.authMW
}
