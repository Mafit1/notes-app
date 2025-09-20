package app

import (
	"context"
	"os"

	"github.com/Mafit1/notes-app/config"
	"github.com/Mafit1/notes-app/internal/api"
	"github.com/Mafit1/notes-app/internal/api/common/middleware"
	"github.com/Mafit1/notes-app/internal/database"
	authRepo "github.com/Mafit1/notes-app/internal/repository/auth"
	notesRepo "github.com/Mafit1/notes-app/internal/repository/notes"
	usersRepo "github.com/Mafit1/notes-app/internal/repository/users"
	authService "github.com/Mafit1/notes-app/internal/service/auth"
	"github.com/Mafit1/notes-app/internal/service/jwtservice"
	notesService "github.com/Mafit1/notes-app/internal/service/notes"
	usersService "github.com/Mafit1/notes-app/internal/service/users"
	"github.com/Mafit1/notes-app/pkg/hasher"
	"github.com/Mafit1/notes-app/pkg/httpserver"
	"github.com/Mafit1/notes-app/pkg/postgres"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	// db
	postgres *postgres.Postgres

	// echo
	echoHandler *echo.Echo

	// repos
	notesRepo notesRepo.Repository
	usersRepo usersRepo.Repository
	authRepo  authRepo.Repository

	// handlers
	postAuthRegisterHandler api.Handler
	postAuthLoginHandler    api.Handler
	postAuthLogoutHandler   api.Handler
	postAuthRefreshHandler  api.Handler

	deleteNoteHandler api.Handler

	getNoteByIDHandler      api.Handler
	getNotesHandler         api.Handler
	getNotesByUserIDHandler api.Handler

	postNoteHandler api.Handler

	putNoteHandler api.Handler

	// services
	notesService notesService.Service
	usersService usersService.Service
	authService  authService.Service
	jwtService   jwtservice.Service

	// infra
	hasher hasher.Hasher

	// middlewares
	authMW *middleware.AuthMW
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

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
