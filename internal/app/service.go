package app

import "github.com/Mafit1/notes-app/internal/service/notes"

func (app *App) NotesService() notes.Service {
	if app.notesService == nil {
		app.notesService = notes.New(app.NotesRepo())
	}
	return app.notesService
}
