package app

import (
	"github.com/Mafit1/notes-app/internal/api"
	deletenote "github.com/Mafit1/notes-app/internal/api/delete_note"
	getnotebyid "github.com/Mafit1/notes-app/internal/api/get_note_by_id"
	getnotes "github.com/Mafit1/notes-app/internal/api/get_notes"
	postnote "github.com/Mafit1/notes-app/internal/api/post_note"
	putnote "github.com/Mafit1/notes-app/internal/api/put_note"
)

func (app *App) DeleteNoteHandler() api.Handler {
	if app.deleteNoteHandler == nil {
		app.deleteNoteHandler = deletenote.New(app.NotesService())
	}
	return app.deleteNoteHandler
}

func (app *App) GetNoteByIDHandler() api.Handler {
	if app.getNoteByIDHandler == nil {
		app.getNoteByIDHandler = getnotebyid.New(app.NotesService())
	}
	return app.getNoteByIDHandler
}

func (app *App) GetNotesHandler() api.Handler {
	if app.getNotesHandler == nil {
		app.getNotesHandler = getnotes.New(app.NotesService())
	}
	return app.getNotesHandler
}

func (app *App) PostNoteHandler() api.Handler {
	if app.postNoteHandler == nil {
		app.postNoteHandler = postnote.New(app.NotesService())
	}
	return app.postNoteHandler
}

func (app *App) PutNoteHandler() api.Handler {
	if app.putNoteHandler == nil {
		app.putNoteHandler = putnote.New(app.NotesService())
	}
	return app.putNoteHandler
}
