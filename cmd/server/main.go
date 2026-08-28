package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"share-platform/internal/config"
	apphttp "share-platform/internal/http"
	"share-platform/internal/repository"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("create database pool", "error", err)
		return
	}
	defer pool.Close()
	server := &http.Server{Addr: cfg.Address, Handler: apphttp.NewRouter(cfg, apphttp.Dependencies{
		Content: repository.NewContent(pool),
		Layout:  repository.NewLayout(pool),
	})}

	go func() {
		slog.Info("server listening", "address", cfg.Address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server stopped", "error", err)
			stop()
		}
	}()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}
