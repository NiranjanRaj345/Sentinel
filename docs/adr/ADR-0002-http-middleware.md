# ADR-0002: HTTP Middleware Pipeline

## Status

Accepted

## Context

As Sentinel grows, HTTP concerns such as panic recovery,
request logging, authentication, request IDs, and rate
limiting should not be duplicated across handlers.

## Decision

Adopt a middleware pipeline using the standard
net/http handler pattern.

Middleware will be transport-specific and live under:

internal/server/middleware

## Consequences

Benefits

- Reusable middleware
- Cleaner handlers
- Easier testing
- Future authentication without handler changes

Trade-offs

- One additional abstraction layer
- Slightly more startup wiring