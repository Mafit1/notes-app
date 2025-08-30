package app

import (
	"github.com/Mafit1/notes-app/pkg/hasher"
)

func (app *App) Hasher() hasher.Hasher {
	if app.hasher == nil {
		app.hasher = hasher.NewBcrypt()
	}
	return app.hasher
}
