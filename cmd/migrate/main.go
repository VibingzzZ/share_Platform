package main

import (
	"context"
	"log"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"

	"share-platform/internal/config"
	"share-platform/internal/db"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool, filepath.Join("db", "migrations")); err != nil {
		log.Fatal(err)
	}
	if err := db.Seed(ctx, pool, filepath.Join("db", "seed")); err != nil {
		log.Fatal(err)
	}
}
