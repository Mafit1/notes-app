package jwtservice

import "github.com/google/uuid"

type GenerateIn struct {
	UserID uuid.UUID
	Email  string
	Role   string
}
