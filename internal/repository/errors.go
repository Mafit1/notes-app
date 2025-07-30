package repository

import "errors"

var (
	ErrDatabase = errors.New("database error")

	ErrNoteNotFound = errors.New("note not found")
)
