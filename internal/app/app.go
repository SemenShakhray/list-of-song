package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SemenShakhray/list-of-song/internal/api/handlers"
	"github.com/SemenShakhray/list-of-song/internal/api/router"
	"github.com/SemenShakhray/list-of-song/internal/config"
	"github.com/SemenShakhray/list-of-song/internal/service"
	"github.com/SemenShakhray/list-of-song/internal/storage/postgres"
	"github.com/SemenShakhray/list-of-song/pkg/logger"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type App struct {
	db     *sql.DB
	server *http.Server
	Sigint chan os.Signal
	cfg    config.Config
}

func (a *App) Run() error {
	log.Printf("Server is start: host - %s, port - %s\n", a.cfg.Server.Host, a.cfg.Server.Port)
	return a.server.ListenAndServe()
}

func (a *App) Stop() error {
	a.db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return a.server.Shutdown(ctx)
}

func NewApp() (*App, error) {

	cfg := config.Load()

	log := logger.NewLogger()
	if log == nil {
		return nil, fmt.Errorf("failed to create logger")
	}
	log.Info("Created logger")

	db, err := postgres.Connect(cfg)
	if err != nil {
		return nil, err
	}
	if db == nil {
		return nil, fmt.Errorf("failed to create DB")
	}

	goose.SetBaseFS(os.DirFS("./migrations"))
	goose.SetDialect("postgres")
	err = goose.Up(db, ".")
	if err != nil {
		return nil, fmt.Errorf("failed up migration: %w", err)
	}
	log.Info("Connected to database and applied migrations", zap.String("config", cfg.DB.User))

	store := postgres.NewStore(db, log)
	if store == nil {
		return nil, fmt.Errorf("failed to create store")
	}

	serv := service.NewService(store)
	if serv == nil {
		return nil, fmt.Errorf("failed to create server")
	}

	handler := handlers.NewHandler(log, serv)

	router := router.NewRouter(&handler)
	if router == nil {
		return nil, fmt.Errorf("failed to create router")
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	server := &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		WriteTimeout: cfg.Server.Timeout,
		ReadTimeout:  cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	log.Info("Created App with server", zap.Any("server", cfg.Server))

	return &App{
		server: server,
		db:     db,
		Sigint: sigint,
		cfg:    cfg,
	}, nil
}
