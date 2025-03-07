package app

import (
	"database/sql"
	"github.com/noueii/rss-gator/internal/config"
	"github.com/noueii/rss-gator/internal/db"
)

import _ "github.com/lib/pq"

type App struct {
	Config *config.Config
	DB     *db.Queries
}

func New() (*App, error) {
	cfg, err := config.Load()

	if err != nil {
		return nil, err
	}

	database, err := sql.Open("postgres", cfg.DbURL)

	if err != nil {
		return nil, err
	}

	dbQueries := db.New(database)

	app := &App{
		Config: cfg,
		DB:     dbQueries,
	}

	return app, nil
}

func (a *App) Close() error {
	return nil
}
