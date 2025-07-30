package notes

import "errors"

var (
	ErrCannotCreateNote  = errors.New("cannot create note")
	ErrCannotGetAllNotes = errors.New("cannot get all notes")
	ErrCannotGetNote     = errors.New("cannot get note")
	ErrNoteNotFound      = errors.New("note not found")
	ErrCannotUpdateNote  = errors.New("cannot update note")
	ErrCannotDeleteNote  = errors.New("cannot delete note")
)
