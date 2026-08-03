# Changelog

All notable changes to Sentinel will be documented in this file.

The format is inspired by Keep a Changelog.

---

## [Unreleased]

### Added

- Recovery engine with policy-driven orchestration
- Guardian hardware integration with REST API
- Observer environmental monitoring integration
- Recovery execution persistence via SQLite
- Guardian and Observer dashboard pages
- Recovery automation actions in rule engine
- Guardian events integrated into event system
- Observer firmware scaffold with PlatformIO

### Changed

- Updated specification with current architecture
- Updated roadmap to reflect completed phases
- Updated backlog with completed items
- README milestone status sections

### Fixed

- Backend JSON field name alignment with frontend types
- API error handling with detailed status and body messages

---

## v0.3.0-dev — Platform Alpha

### Added

- Complete monitoring subsystem
- Event-driven architecture
- Rules engine with condition/action model
- Automation engine
- Recovery engine with retry and delay support
- Guardian integration (power/reset/status)
- Observer integration (temperature/humidity/status)
- Dashboard pages for all subsystems
- SQLite persistence for events, rules, automation, recovery
- WebSocket streaming
- Authentication middleware
- Next.js dashboard with React Query

---

## v0.2.0 — Core Foundation

### Added
- HTTP server
- Middleware pipeline
- Health endpoint
- Metrics endpoint
- Version package
- Metadata collection
- Graceful shutdown
- HTTP handler tests

### Changed
- Centralized version information
- Refactored HTTP handlers
- Introduced application lifecycle

### Removed
- Legacy internal/server/handlers.go

---

## [0.8.0]

### Added

- Development automation scripts
- Quality gate workflow

---

## [0.7.0]

### Added

- `/system` endpoint
- System package

---

## [0.6.0]

### Added

- `/health` endpoint

---

## [0.5.0]

### Added

- HTTP server

---

## [0.4.0]

### Added

- Sentinel specification

---

## [0.3.0]

### Added

- Application lifecycle

---

## [0.2.0]

### Added

- Executable Node Agent

---

## [0.1.0]

### Added

- Repository initialization
