package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Seed executes repeatable, idempotent demo data scripts in lexical order.
func Seed(ctx context.Context, pool *pgxpool.Pool, seedDir string) error {
	seeds, err := sqlFiles(seedDir)
	if err != nil {
		return fmt.Errorf("load seed files: %w", err)
	}
	for _, seed := range seeds {
		if _, err := pool.Exec(ctx, seed.sql); err != nil {
			return fmt.Errorf("apply seed %q: %w", seed.name, err)
		}
	}
	return nil
}
