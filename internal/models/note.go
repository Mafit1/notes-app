package models

import "github.com/google/uuid"

type Note struct {
	ID      int64
	Title   string
	Content string
	UserID  uuid.UUID
}
