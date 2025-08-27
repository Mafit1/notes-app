package app

import (
	"github.com/Mafit1/notes-app/pkg/hasher"
	"github.com/Mafit1/notes-app/pkg/uservalidator"
)

func (app *App) PasswordHasher() hasher.Hasher {
	if app.passwordHasher == nil {
		app.passwordHasher = hasher.NewBcrypt()
	}
	return app.passwordHasher
}

func (app *App) UserValidator() uservalidator.UserValidator {
	if app.userValidator == nil {
		app.userValidator = uservalidator.New()
	}
	return app.userValidator
}
