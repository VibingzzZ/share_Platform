package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Migrate applies every SQL migration not already recorded in schema_migrations.
func Migrate(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	if _, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			name TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("create migration ledger: %w", err)
	}

	migrations, err := sqlFiles(migrationsDir)
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}
	for _, migration := range migrations {
		var applied bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE name = $1)`, migration.name).Scan(&applied); err != nil {
			return fmt.Errorf("check migration %q: %w", migration.name, err)
		}
		if applied {
			continue
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin migration %q: %w", migration.name, err)
		}
		if _, err := tx.Exec(ctx, migration.sql); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("apply migration %q: %w", migration.name, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1)`, migration.name); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("record migration %q: %w", migration.name, err)
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %q: %w", migration.name, err)
		}
	}
	return nil
}

type sqlFile struct {
	name string
	sql  string
}

func sqlFiles(dir string) ([]sqlFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]sqlFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		contents, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		files = append(files, sqlFile{name: entry.Name(), sql: string(contents)})
	}
	sort.Slice(files, func(i, j int) bool { return strings.Compare(files[i].name, files[j].name) < 0 })
	return files, nil
}
