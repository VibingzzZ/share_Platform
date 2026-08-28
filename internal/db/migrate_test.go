package db

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMigrateAndSeed(t *testing.T) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is unset")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()
	if err := p.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	root := filepath.Join("..", "..")
	if err := Migrate(ctx, p, filepath.Join(root, "db", "migrations")); err != nil {
		t.Fatal(err)
	}
	if err := Seed(ctx, p, filepath.Join(root, "db", "seed")); err != nil {
		t.Fatal(err)
	}

	var count int
	if err := p.QueryRow(ctx, `SELECT count(*) FROM users WHERE email = $1`, "member@example.com").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected one demo member, got %d", count)
	}
	if err := Seed(ctx, p, filepath.Join(root, "db", "seed")); err != nil {
		t.Fatal(err)
	}
	var total int
	if err := p.QueryRow(ctx, `SELECT count(*) FROM users WHERE email = $1`, "member@example.com").Scan(&total); err != nil {
		t.Fatal(err)
	}
	if total != 1 {
		t.Fatalf("seed is not idempotent: got %d demo members", total)
	}
}
