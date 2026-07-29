# ADR-0001

## Title

Layered Architecture

## Status

Accepted

## Decision

Sentinel adopts a layered architecture.

Application

↓

Server

↓

Service

↓

System

## Consequences

Positive

- Clear responsibilities
- Testability
- Loose coupling

Negative

- Slightly more boilerplate