package jwtservice

import "github.com/google/uuid"

type GenerateIn struct {
	UserID uuid.UUID
	Email  string
	Role   string
}

type GenerateOut struct {
	AccessToken  string
	RefreshToken string
}
