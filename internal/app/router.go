package app

import "github.com/labstack/echo/v4"

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()

	app.configureRouter(handler)

	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	notesGroup := handler.Group("/notes")
	{
		notesGroup.GET("", app.GetNotesHandler().Handle)
		notesGroup.GET("/:id", app.GetNoteByIDHandler().Handle)
		notesGroup.POST("/:id", app.PostNoteHandler().Handle)
		notesGroup.DELETE("/:id", app.DeleteNoteHandler().Handle)
	}
}
