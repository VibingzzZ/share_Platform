package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"share-platform/internal/model"
)

type LayoutRepository struct {
	pool *pgxpool.Pool
}

func NewLayout(pool *pgxpool.Pool) *LayoutRepository {
	return &LayoutRepository{pool: pool}
}

func (r *LayoutRepository) LoadLayout(ctx context.Context, userID string) (model.Layout, error) {
	var raw []byte
	err := r.pool.QueryRow(ctx, `SELECT layout FROM user_layouts WHERE user_id = $1`, userID).Scan(&raw)
	if err == pgx.ErrNoRows {
		return model.DefaultLayout(), nil
	}
	if err != nil {
		return model.Layout{}, fmt.Errorf("load layout: %w", err)
	}
	var layout model.Layout
	if err := json.Unmarshal(raw, &layout); err != nil {
		return model.Layout{}, fmt.Errorf("decode layout: %w", err)
	}
	if err := layout.Validate(); err != nil {
		return model.Layout{}, fmt.Errorf("validate stored layout: %w", err)
	}
	return layout, nil
}

func (r *LayoutRepository) SaveLayout(ctx context.Context, userID string, layout model.Layout) error {
	if err := layout.Validate(); err != nil {
		return err
	}
	raw, err := json.Marshal(layout)
	if err != nil {
		return fmt.Errorf("encode layout: %w", err)
	}
	_, err = r.pool.Exec(ctx, `INSERT INTO user_layouts (user_id, layout) VALUES ($1, $2::jsonb)
		ON CONFLICT (user_id) DO UPDATE SET layout = EXCLUDED.layout, updated_at = now()`, userID, raw)
	if err != nil {
		return fmt.Errorf("save layout: %w", err)
	}
	return nil
}
