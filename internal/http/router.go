package http

import (
	"encoding/json"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
	"share-platform/internal/config"
)

// Dependencies groups the storage adapters used by HTTP handlers.
type Dependencies struct {
	Content ContentStore
	Layout  LayoutStore
}

func NewRouter(cfg config.Config, deps Dependencies) nethttp.Handler {
	_ = cfg
	r := chi.NewRouter()
	r.Use(identityMiddleware)
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler)
		r.Get("/dashboard", dashboardHandler(deps.Content))
		r.Get("/resources", resourcesHandler(deps.Content))
		r.Get("/posts", postsHandler(deps.Content))
		r.Get("/ai-products", aiProductsHandler(deps.Content))
		r.Group(func(r chi.Router) {
			r.Use(requireIdentity)
			r.Get("/layout", getLayoutHandler(deps.Layout))
			r.Put("/layout", putLayoutHandler(deps.Layout))
		})
		r.Route("/admin", func(r chi.Router) {
			r.Use(requireAdmin)
			r.Post("/resources", createResourceHandler(deps.Content))
			r.Put("/resources/{id}", updateResourceHandler(deps.Content))
			r.Delete("/resources/{id}", deleteResourceHandler(deps.Content))
			r.Post("/posts", createPostHandler(deps.Content))
			r.Put("/posts/{id}", updatePostHandler(deps.Content))
			r.Delete("/posts/{id}", deletePostHandler(deps.Content))
			r.Post("/ai-products", createAIProductHandler(deps.Content))
			r.Put("/ai-products/{id}", updateAIProductHandler(deps.Content))
			r.Delete("/ai-products/{id}", deleteAIProductHandler(deps.Content))
		})
	})
	return r
}

func healthHandler(w nethttp.ResponseWriter, _ *nethttp.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
