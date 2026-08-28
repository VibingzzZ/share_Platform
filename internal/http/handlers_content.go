package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	nethttp "net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"share-platform/internal/model"
)

type ContentStore interface {
	ListResources(context.Context, model.ListFilter) (model.Page[model.Resource], error)
	ListPosts(context.Context, model.ListFilter) (model.Page[model.Post], error)
	ListAIProducts(context.Context, model.ListFilter) (model.Page[model.AIProduct], error)
	Dashboard(context.Context) (model.Dashboard, error)
	CreateResource(context.Context, model.Resource) (model.Resource, error)
	UpdateResource(context.Context, string, model.Resource) (model.Resource, error)
	DeleteResource(context.Context, string) error
	CreatePost(context.Context, model.Post) (model.Post, error)
	UpdatePost(context.Context, string, model.Post) (model.Post, error)
	DeletePost(context.Context, string) error
	CreateAIProduct(context.Context, model.AIProduct) (model.AIProduct, error)
	UpdateAIProduct(context.Context, string, model.AIProduct) (model.AIProduct, error)
	DeleteAIProduct(context.Context, string) error
}

func dashboardHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		dashboard, err := store.Dashboard(r.Context())
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, dashboard)
	}
}

func resourcesHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		filter, err := listFilter(r)
		if err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		page, err := store.ListResources(r.Context(), filter)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, page)
	}
}

func postsHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		filter, err := listFilter(r)
		if err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		page, err := store.ListPosts(r.Context(), filter)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, page)
	}
}

func aiProductsHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		filter, err := listFilter(r)
		if err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		page, err := store.ListAIProducts(r.Context(), filter)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, page)
	}
}

func listFilter(r *nethttp.Request) (model.ListFilter, error) {
	filter := model.ListFilter{Page: 1, PageSize: 20, Type: r.URL.Query().Get("type"), Tag: r.URL.Query().Get("tag")}
	for key, target := range map[string]*int{"page": &filter.Page, "pageSize": &filter.PageSize} {
		if raw := r.URL.Query().Get(key); raw != "" {
			value, err := strconv.Atoi(raw)
			if err != nil || value < 1 {
				return model.ListFilter{}, fmt.Errorf("%s must be a positive integer", key)
			}
			*target = value
		}
	}
	if filter.PageSize > 100 {
		return model.ListFilter{}, errors.New("pageSize must not exceed 100")
	}
	return filter, nil
}

func createResourceHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var resource model.Resource
		if !decodeJSON(w, r, &resource) {
			return
		}
		resource.AuthorID = currentIdentity(r.Context()).UserID
		if err := validateResource(resource); err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		created, err := store.CreateResource(r.Context(), resource)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusCreated, created)
	}
}

func updateResourceHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var resource model.Resource
		if !decodeJSON(w, r, &resource) {
			return
		}
		if err := validateResource(resource); err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		updated, err := store.UpdateResource(r.Context(), chi.URLParam(r, "id"), resource)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, updated)
	}
}

func deleteResourceHandler(store ContentStore) nethttp.HandlerFunc {
	return deleteHandler(store, func(ctx context.Context, id string) error { return store.DeleteResource(ctx, id) })
}

func createPostHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var post model.Post
		if !decodeJSON(w, r, &post) {
			return
		}
		post.AuthorID = currentIdentity(r.Context()).UserID
		if err := validatePost(post); err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		created, err := store.CreatePost(r.Context(), post)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusCreated, created)
	}
}

func updatePostHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var post model.Post
		if !decodeJSON(w, r, &post) {
			return
		}
		if err := validatePost(post); err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		updated, err := store.UpdatePost(r.Context(), chi.URLParam(r, "id"), post)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, updated)
	}
}

func deletePostHandler(store ContentStore) nethttp.HandlerFunc {
	return deleteHandler(store, func(ctx context.Context, id string) error { return store.DeletePost(ctx, id) })
}

func createAIProductHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var product model.AIProduct
		if !decodeJSON(w, r, &product) {
			return
		}
		if err := validateAIProduct(product); err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		created, err := store.CreateAIProduct(r.Context(), product)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusCreated, created)
	}
}

func updateAIProductHandler(store ContentStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var product model.AIProduct
		if !decodeJSON(w, r, &product) {
			return
		}
		if err := validateAIProduct(product); err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		updated, err := store.UpdateAIProduct(r.Context(), chi.URLParam(r, "id"), product)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, updated)
	}
}

func deleteAIProductHandler(store ContentStore) nethttp.HandlerFunc {
	return deleteHandler(store, func(ctx context.Context, id string) error { return store.DeleteAIProduct(ctx, id) })
}

func deleteHandler(store ContentStore, deleteItem func(context.Context, string) error) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "content service is unavailable")
			return
		}
		if err := deleteItem(r.Context(), chi.URLParam(r, "id")); err != nil {
			writeStoreError(w, err)
			return
		}
		w.WriteHeader(nethttp.StatusNoContent)
	}
}

func decodeJSON(w nethttp.ResponseWriter, r *nethttp.Request, target any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		writeError(w, nethttp.StatusBadRequest, "bad_request", "invalid JSON body")
		return false
	}
	return true
}

func validateResource(resource model.Resource) error {
	if resource.Title == "" || resource.Type == "" || resource.Status == "" {
		return errors.New("title, type, and status are required")
	}
	if !model.ValidStatus(resource.Status) {
		return errors.New("invalid status")
	}
	if resource.URL == nil && resource.FilePath == nil {
		return errors.New("url or filePath is required")
	}
	return nil
}

func validatePost(post model.Post) error {
	if post.Title == "" || post.Body == "" || post.Status == "" {
		return errors.New("title, body, and status are required")
	}
	if !model.ValidStatus(post.Status) {
		return errors.New("invalid status")
	}
	return nil
}

func validateAIProduct(product model.AIProduct) error {
	if product.Name == "" || product.URL == "" || product.Status == "" {
		return errors.New("name, url, and status are required")
	}
	if !model.ValidStatus(product.Status) {
		return errors.New("invalid status")
	}
	return nil
}

func writeStoreError(w nethttp.ResponseWriter, err error) {
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, nethttp.StatusNotFound, "not_found", "content was not found")
		return
	}
	writeError(w, nethttp.StatusInternalServerError, "internal_error", "request could not be completed")
}
