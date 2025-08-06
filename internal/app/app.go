package app

import (
	"context"
	"os"

	"github.com/Mafit1/notes-app/config"
	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/database"
	notesRepo "github.com/Mafit1/notes-app/internal/repository/notes"
	notesService "github.com/Mafit1/notes-app/internal/service/notes"
	"github.com/Mafit1/notes-app/pkg/httpserver"
	"github.com/Mafit1/notes-app/pkg/postgres"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	postgres *postgres.Postgres

	echoHandler *echo.Echo

	notesRepo notesRepo.Repository

	deleteNoteHandler api.Handler

	getNoteByIDHandler api.Handler
	getNotesHandler    api.Handler

	postNoteHandler api.Handler

	notesService notesService.Service
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	return &App{cfg: cfg}
}

func (app *App) Start() {
	log.Info("Connecting to PostgreSQL...")
	postgres, err := postgres.New(app.cfg.Postgres.URL, postgres.ConnAttempts(5))
	if err != nil {
		log.Fatalf("app - Start - Postgres failed: %v", err)
	}

	app.postgres = postgres

	defer postgres.Close()

	err = database.RunMigrations(context.Background(), app.postgres.Pool)
	if err != nil {
		log.Fatalf("app - Start - Migrations failed: %v", err)
	}

	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()

	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			log.Errorf("HTTP server shutdown error: %v", err)
		}
	}()

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	log.Info("Shutting down...")
}
