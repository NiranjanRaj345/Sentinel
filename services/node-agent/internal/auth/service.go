package auth

import (
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
)

func FromConfig(cfg config.AuthConfig) (*TokenStore, error) {
	tokens := make([]Token, 0, len(cfg.Tokens))
	for _, tokenCfg := range cfg.Tokens {
		permissions := make([]Permission, 0, len(tokenCfg.Permissions))
		for _, perm := range tokenCfg.Permissions {
			permissions = append(permissions, Permission(perm))
		}

		tokens = append(tokens, Token{
			ID:          tokenCfg.Name,
			Name:        tokenCfg.Name,
			Value:       tokenCfg.Value,
			Permissions: permissions,
			ExpiresAt:   tokenCfg.ExpiresAt,
		})
	}

	return NewTokenStore(tokens), nil
}

func IsExpired(token Token) bool {
	if token.ExpiresAt == nil {
		return false
	}
	return time.Now().UTC().After(*token.ExpiresAt)
}
