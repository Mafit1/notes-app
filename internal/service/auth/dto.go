package auth

import "github.com/Mafit1/notes-app/internal/models"

type RegisterIn struct {
	name     string
	email    string
	password string
}

type LoginIn struct {
	email    string
	password string
}

type RegisterOut struct {
	user      *models.User
	authToken string
}

type LoginOut struct {
	user      *models.User
	authToken string
}
