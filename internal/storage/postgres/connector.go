package postgres

import (
	"database/sql"
	"fmt"

	"github.com/SemenShakhray/list-of-song/internal/config"
	"github.com/SemenShakhray/list-of-song/internal/storage"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Store struct {
	DB  *sql.DB
	Log *zap.Logger
}

func NewStore(db *sql.DB, log *zap.Logger) storage.Storer {
	return &Store{
		DB:  db,
		Log: log,
	}
}

func Connect(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}
