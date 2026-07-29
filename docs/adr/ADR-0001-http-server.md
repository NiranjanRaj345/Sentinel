# ADR-0001: HTTP Server Implementation

## Status

Accepted

## Context

Sentinel requires an HTTP interface for future APIs.

## Decision

Use net/http with an explicit http.Server instance.

Do not use the global DefaultServeMux.

## Consequences

Advantages

- Graceful shutdown support
- Explicit routing
- Better testability
- Future TLS support

Disadvantages

- Slightly more boilerplate