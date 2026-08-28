package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"share-platform/internal/config"
)

// Dependencies groups services used by HTTP handlers. It is intentionally
// small at the scaffold stage and can grow without changing NewRouter.
type Dependencies struct{}

func NewRouter(cfg config.Config, deps Dependencies) nethttp.Handler {
	_ = cfg
	_ = deps
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler)
	})
	return r
}

func healthHandler(w nethttp.ResponseWriter, _ *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
