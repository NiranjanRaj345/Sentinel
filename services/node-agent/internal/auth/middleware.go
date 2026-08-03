package auth

import (
	"net/http"
	"strings"
)

func Authenticate(store *TokenStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenValue := extractToken(r)
			if tokenValue == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			token, ok := store.Resolve(tokenValue)
			if !ok {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := ContextWithToken(r.Context(), token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractToken(r *http.Request) string {
	header := r.Header.Get("Authorization")
	if header != "" {
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	token := r.URL.Query().Get("token")
	if token != "" {
		return token
	}

	return ""
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
