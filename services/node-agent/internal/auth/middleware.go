package auth

import (
	"net/http"
	"strings"
)

func Authenticate(store *TokenStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			token, ok := store.Resolve(parts[1])
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := ContextWithToken(r.Context(), token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Authorize(permission Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := TokenFromContext(r.Context())
			if !ok {
				http.Error(w, "missing token context", http.StatusInternalServerError)
				return
			}

			if !HasPermission(token, permission) {
				http.Error(w, "insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
