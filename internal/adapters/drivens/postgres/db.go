package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/anderson89marques/roxb3/internal/infra/config"
	_ "github.com/lib/pq"
)

func NewPostgresDB(conf *config.Config) (*sql.DB, error) {
	dsn := conf.DatabaseURL()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error to connect to database: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("file to ping database connection: %w", err)
	}
	return db, nil
}
