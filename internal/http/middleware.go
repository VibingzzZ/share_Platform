package http

import (
	"context"
	"encoding/json"
	nethttp "net/http"
)

type identity struct {
	UserID string
	Role   string
}

type identityContextKey struct{}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func identityMiddleware(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		role := r.Header.Get("X-User-Role")
		if role == "" {
			role = "member"
		}
		ctx := context.WithValue(r.Context(), identityContextKey{}, identity{UserID: r.Header.Get("X-User-ID"), Role: role})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func currentIdentity(ctx context.Context) identity {
	value, _ := ctx.Value(identityContextKey{}).(identity)
	return value
}

func requireIdentity(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if currentIdentity(r.Context()).UserID == "" {
			writeError(w, nethttp.StatusUnauthorized, "unauthorized", "X-User-ID header is required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func requireAdmin(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		principal := currentIdentity(r.Context())
		if principal.UserID == "" {
			writeError(w, nethttp.StatusUnauthorized, "unauthorized", "X-User-ID header is required")
			return
		}
		if principal.Role != "admin" {
			writeError(w, nethttp.StatusForbidden, "forbidden", "administrator role is required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w nethttp.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w nethttp.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Code: code, Message: message})
}
