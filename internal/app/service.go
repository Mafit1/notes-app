package app

import (
	"github.com/Mafit1/notes-app/internal/service/auth"
	"github.com/Mafit1/notes-app/internal/service/jwtservice"
	"github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/Mafit1/notes-app/internal/service/users"
)

func (app *App) NotesService() notes.Service {
	if app.notesService == nil {
		app.notesService = notes.New(app.NotesRepo())
	}
	return app.notesService
}

func (app *App) UsersService() users.Service {
	if app.usersService == nil {
		app.usersService = users.New(app.UsersRepo())
	}
	return app.usersService
}

func (app *App) AuthService() auth.Service {
	if app.authService == nil {
		app.authService = auth.New(app.UsersService(), app.JwtService(), app.UserValidator(), app.PasswordHasher())
	}
	return app.authService
}

func (app *App) JwtService() jwtservice.Service {
	if app.jwtservice == nil {
		app.jwtservice = jwtservice.New(app.cfg.Auth.JWTAccessSecretKey, app.cfg.Auth.AccessTokenTTL)
	}
	return app.jwtservice
}
