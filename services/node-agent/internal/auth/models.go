package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type Permission string

const (
	PermissionRead      Permission = "read"
	PermissionOperate   Permission = "operate"
	PermissionConfigure Permission = "configure"
)

type Token struct {
	ID          string
	Name        string
	Value       string
	Permissions []Permission
	ExpiresAt   *time.Time
}

type TokenStore struct {
	tokens map[string]Token
	mu     sync.RWMutex
}

func NewTokenStore(tokens []Token) *TokenStore {
	store := &TokenStore{tokens: make(map[string]Token)}
	for _, token := range tokens {
		store.tokens[token.Value] = token
	}
	return store
}

func (s *TokenStore) Resolve(value string) (Token, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	token, ok := s.tokens[value]
	return token, ok
}

func (s *TokenStore) HasPermission(token Token, permission Permission) bool {
	for _, p := range token.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func HasPermission(token Token, permission Permission) bool {
	for _, p := range token.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

type contextKey string

const TokenContextKey contextKey = "token"

func ContextWithToken(ctx context.Context, token Token) context.Context {
	return context.WithValue(ctx, TokenContextKey, token)
}

func TokenFromContext(ctx context.Context) (Token, bool) {
	token, ok := ctx.Value(TokenContextKey).(Token)
	return token, ok
}

func GenerateTokenValue() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
