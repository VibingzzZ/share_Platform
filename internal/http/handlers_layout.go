package http

import (
	"context"
	"errors"
	nethttp "net/http"

	"share-platform/internal/model"
)

type LayoutStore interface {
	LoadLayout(context.Context, string) (model.Layout, error)
	SaveLayout(context.Context, string, model.Layout) error
}

func getLayoutHandler(store LayoutStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "layout service is unavailable")
			return
		}
		layout, err := store.LoadLayout(r.Context(), currentIdentity(r.Context()).UserID)
		if err != nil {
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, layout)
	}
}

func putLayoutHandler(store LayoutStore) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {
		var layout model.Layout
		if !decodeJSON(w, r, &layout) {
			return
		}
		if err := layout.Validate(); err != nil {
			writeError(w, nethttp.StatusBadRequest, "bad_request", err.Error())
			return
		}
		if store == nil {
			writeError(w, nethttp.StatusServiceUnavailable, "service_unavailable", "layout service is unavailable")
			return
		}
		if err := store.SaveLayout(r.Context(), currentIdentity(r.Context()).UserID, layout); err != nil {
			if errors.Is(err, context.Canceled) {
				writeError(w, nethttp.StatusRequestTimeout, "request_cancelled", "request was cancelled")
				return
			}
			writeStoreError(w, err)
			return
		}
		writeJSON(w, nethttp.StatusOK, layout)
	}
}
