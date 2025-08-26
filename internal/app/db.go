package app

import (
	"github.com/Mafit1/notes-app/internal/repository/notes"
	"github.com/Mafit1/notes-app/internal/repository/users"
	"github.com/Mafit1/notes-app/pkg/postgres"
)

func (app *App) Postgres() *postgres.Postgres {
	return app.postgres
}

func (app *App) NotesRepo() notes.Repository {
	if app.notesRepo == nil {
		app.notesRepo = notes.New(app.Postgres())
	}
	return app.notesRepo
}

func (app *App) UsersRepo() users.Repository {
	if app.usersRepo == nil {
		app.usersRepo = users.New(app.Postgres())
	}
	return app.usersRepo
}
