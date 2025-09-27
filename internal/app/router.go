package app

import (
	"github.com/Mafit1/notes-app/pkg/validator"
	"github.com/labstack/echo/v4"
)

func (app *App) EchoHandler() *echo.Echo {
	if app.echoHandler != nil {
		return app.echoHandler
	}

	handler := echo.New()
	handler.Validator = validator.New()

	app.configureRouter(handler)
	app.echoHandler = handler

	return app.echoHandler
}

func (app *App) configureRouter(handler *echo.Echo) {
	// Metrics
	handler.Use(app.MetricsMW().MetricsMiddleware)

	// Auth
	handler.POST("/register", app.PostAuthRegisterHandler().Handle)
	handler.POST("/login", app.PostAuthLoginHandler().Handle)
	handler.POST("/logout", app.PostAuthLogoutHandler().Handle)
	handler.POST("/refresh", app.PostAuthRefreshHandler().Handle)

	notesGroup := handler.Group("/notes", app.AuthMW().Authenticate())
	{
		notesGroup.GET("", app.GetNotesByUserIDHandler().Handle)
		notesGroup.GET("/:id", app.GetNoteByIDHandler().Handle)
		notesGroup.POST("", app.PostNoteHandler().Handle)
		notesGroup.DELETE("/:id", app.DeleteNoteHandler().Handle)
		notesGroup.PUT("/:id", app.PutNoteHandler().Handle)
	}
}
