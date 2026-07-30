package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenStore_Resolve(t *testing.T) {
	store := NewTokenStore([]Token{
		{Name: "test", Value: "secret", Permissions: []Permission{PermissionRead}},
	})

	token, ok := store.Resolve("secret")
	if !ok {
		t.Fatal("expected token to be found")
	}
	if token.Name != "test" {
		t.Fatalf("expected name test, got %s", token.Name)
	}
}

func TestTokenStore_Resolve_NotFound(t *testing.T) {
	store := NewTokenStore(nil)
	_, ok := store.Resolve("missing")
	if ok {
		t.Fatal("expected token not found")
	}
}

func TestAuthenticate_MissingHeader_Returns401(t *testing.T) {
	store := NewTokenStore(nil)
	handler := Authenticate(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_InvalidHeader_Returns401(t *testing.T) {
	store := NewTokenStore(nil)
	handler := Authenticate(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Invalid")
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_InvalidToken_Returns401(t *testing.T) {
	store := NewTokenStore([]Token{{Name: "test", Value: "secret", Permissions: []Permission{PermissionRead}}})
	handler := Authenticate(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuthenticate_ValidToken_PassesThrough(t *testing.T) {
	store := NewTokenStore([]Token{{Name: "test", Value: "secret", Permissions: []Permission{PermissionRead}}})
	handler := Authenticate(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAuthorize_InsufficientPermission_Returns403(t *testing.T) {
	store := NewTokenStore([]Token{{Name: "test", Value: "secret", Permissions: []Permission{PermissionRead}}})
	authMiddleware := Authenticate(store)
	authz := Authorize(PermissionOperate)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := authMiddleware(authz(inner))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected %d, got %d", http.StatusForbidden, rr.Code)
	}
}

func TestContextToken_RoundTrip(t *testing.T) {
	token := Token{Name: "test", Value: "secret"}
	ctx := ContextWithToken(context.Background(), token)

	got, ok := TokenFromContext(ctx)
	if !ok {
		t.Fatal("expected token in context")
	}
	if got.Name != "test" {
		t.Fatalf("expected name test, got %s", got.Name)
	}
}

func TestIsExpired_NilExpiry_ReturnsFalse(t *testing.T) {
	token := Token{Name: "test"}
	if IsExpired(token) {
		t.Fatal("expected non-expired token")
	}
}

func TestIsExpired_PastExpiry_ReturnsTrue(t *testing.T) {
	past := time.Now().UTC().Add(-time.Hour)
	token := Token{Name: "test", ExpiresAt: &past}
	if !IsExpired(token) {
		t.Fatal("expected expired token")
	}
}

func TestTokenStore_Empty_NoTokens(t *testing.T) {
	store := NewTokenStore(nil)
	if !store.Empty() {
		t.Fatal("expected empty store")
	}
}

func TestTokenStore_Empty_WithTokens(t *testing.T) {
	store := NewTokenStore([]Token{{Name: "test", Value: "secret"}})
	if store.Empty() {
		t.Fatal("expected non-empty store")
	}
}
