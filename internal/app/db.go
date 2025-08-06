package app

import (
	"github.com/Mafit1/notes-app/internal/repository/notes"
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
