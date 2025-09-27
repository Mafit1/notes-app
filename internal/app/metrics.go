package app

import "github.com/Mafit1/notes-app/internal/metrics"

func (app *App) NotesMetrics() metrics.NotesMetrics {
	if app.notesMetrics == nil {
		app.notesMetrics = metrics.NewNotesMetrics()
	}
	return app.notesMetrics
}
