package app

import (
	"github.com/Mafit1/notes-app/pkg/hasher"
	"github.com/Mafit1/notes-app/pkg/uservalidator"
)

func (app *App) Hasher() hasher.Hasher {
	if app.hasher == nil {
		app.hasher = hasher.NewBcrypt()
	}
	return app.hasher
}

func (app *App) UserValidator() uservalidator.UserValidator {
	if app.userValidator == nil {
		app.userValidator = uservalidator.New()
	}
	return app.userValidator
}
