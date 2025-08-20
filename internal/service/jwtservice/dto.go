package jwtservice

import "github.com/google/uuid"

type RegisterIn struct {
	UserID uuid.UUID
	Email  string
	Role   string
}
