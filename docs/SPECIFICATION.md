# 🛡️ Sentinel Specification

**Version:** 1.0.0-alpha.1  
**Status:** Active Development

---

# 1. Purpose

Sentinel is a self-hosted personal infrastructure platform designed to monitor, manage, protect, and recover computing devices from anywhere in the world.

The objective of Sentinel is to provide a reliable, secure, extensible, and maintainable platform that remains useful for many years.

---

# 2. Vision

Sentinel begins with a single desktop computer.

It evolves into a unified platform capable of managing:

- Desktop computers
- Linux servers
- NAS systems
- Raspberry Pi devices
- ESP32 hardware
- Home laboratory infrastructure
- Future AI infrastructure

Today, Sentinel is a **cyber-physical infrastructure platform** with software observability, automation, recovery orchestration, and dedicated hardware subsystems for remote power control and environmental awareness.

---

# 3. Engineering Principles

Sentinel follows these principles.

## 3.1 Reliability First

Reliability is always preferred over additional features.

## 3.2 Simplicity

Prefer the simplest solution that correctly solves the problem.

## 3.3 Extensibility

Every component should be designed so future functionality can be added with minimal modification.

## 3.4 Single Responsibility

Every package, file and component should have exactly one responsibility.

## 3.5 Explicit Design

Architecture decisions must be documented before implementation.

## 3.6 Documentation

Important decisions must be recorded.

Documentation is considered part of the source code.

---

# 4. Repository Structure

```
Sentinel/

apps/
  dashboard/          # Next.js dashboard
services/
  node-agent/         # Core backend service
    internal/
      alert/          # Alert evaluation engine
      application/    # Application lifecycle
      auth/           # Token-based authentication
      automation/     # Automation engine + SQLite storage
      config/         # YAML configuration
      dashboard/      # Dashboard aggregation + WebSocket hub
      events/         # Event service + SQLite storage
      guardian/       # Guardian hardware client + service
      logger/         # Structured logging
      metrics/        # Metrics collection + snapshot history
      node/           # Node identity + Linux provider
      observer/       # Observer environmental client + service
      operations/     # Operations service + Linux provider
      recovery/       # Recovery engine + SQLite storage
      resources/      # Resources service + providers
      rules/          # Rule engine + SQLite storage
      scheduler/      # Periodic metrics scheduler
      server/         # HTTP server, routes, middleware
      service/        # Shared service types
      services/       # Services service + providers
      storage/sqlite/ # Primary SQLite store
      stream/         # WebSocket hub
    cmd/
      sentinel-agent/ # Entrypoint
firmware/
  guardian/           # ESP32 Guardian firmware
  observer/           # ESP32 Observer firmware
docs/                 # Specification, ADRs, changelog
```

---

# 5. Platform Architecture

```text
                  Dashboard
                       │
                Authentication
                       │
        REST ──────────┼───────── WebSocket
                       │
     ┌─────────────────┼──────────────────┐
     ▼                 ▼                  ▼
 Resources         Automation          Guardian
     │                 │                  │
     ▼                 ▼                  ▼
 Rule Engine ─────► Operations      ESP32 Client
     │                                    │
     ▼                                    ▼
 Event Service                     Wi-Fi Network
     │                                    │
     ▼                                    ▼
 SQLite                          Relay GPIO
                                          │
                                          ▼
                                  Desktop Power Button
```

---

# 6. Core Subsystems

## 6.1 Node Agent

Runs on every managed machine.

Responsibilities:

- System metrics collection
- Health reporting
- Service and resource management
- Event publishing
- Rule evaluation
- Automation execution
- Recovery orchestration
- Guardian and Observer client integration

---

## 6.2 Dashboard

Primary user interface.

Responsibilities:

- Overview cards
- Live metrics streaming
- History visualization
- Event timeline
- Rule management
- Service and resource control
- Guardian controls
- Recovery execution
- Observer environmental data

---

## 6.3 Events

Event-driven communication layer.

Responsibilities:

- Event publishing
- Event persistence
- Event history
- Source segregation: operations, alerts, scheduler, resources, guardian

---

## 6.4 Rules

Condition/action rule engine.

Responsibilities:

- Rule storage
- Rule evaluation against events
- Action dispatch: notify, execute, guardian_power, guardian_reset

---

## 6.5 Automation

Action dispatcher.

Responsibilities:

- Receive matched rules
- Execute operations
- Invoke Guardian actions
- Publish automation executions

---

## 6.6 Recovery

Recovery orchestration engine.

Responsibilities:

- Execute recovery policies
- Step sequencing with delay and retries
- Guardian integration
- Execution persistence
- Status and recent execution API

---

## 6.7 Guardian

Hardware recovery controller.

Responsibilities:

- Power button relay control
- Reset button relay control
- Power LED monitoring
- Status and health API
- Event publishing for power/reset actions

Sentinel does not know about GPIO pins. It only calls:

```go
guardian.Power()
guardian.Reset()
guardian.Status()
```

---

## 6.8 Observer

Environmental monitoring subsystem.

Responsibilities:

- Temperature and humidity sensing
- Local display support
- Status and environment API
- Event publishing for environmental thresholds

Sentinel only sees:

```go
observer.Status()
observer.Environment()
```

---

# 7. REST API

## Agent

| Method | Path | Purpose |
| ------ | ---- | ------- |
| GET | `/health` | Service health |
| GET | `/metrics` | Latest metrics snapshot |
| GET | `/history/latest` | Latest history entry |
| GET | `/history` | History range |
| GET | `/stream` | WebSocket metrics stream |
| GET | `/dashboard/overview` | Dashboard overview |
| GET | `/dashboard/stream` | Dashboard WebSocket |
| GET | `/dashboard/history` | Dashboard history |
| GET | `/dashboard/capabilities` | Node capabilities |
| GET | `/events/recent` | Recent events |
| GET | `/rules` | Rules |
| POST | `/operations` | Execute operation |
| POST | `/services` | Service action |
| GET | `/resources` | Resources |
| GET | `/automation/executions` | Automation executions |
| GET | `/guardian/status` | Guardian status |
| POST | `/guardian/power` | Power relay action |
| POST | `/guardian/reset` | Reset relay action |
| GET | `/observer/status` | Observer status |
| GET | `/observer/environment` | Environment readings |
| POST | `/recovery/execute` | Execute recovery policy |
| GET | `/recovery/recent` | Recent recovery executions |

---

# 8. Dashboard Pages

| Path | Purpose |
| ---- | ------- |
| `/` | Overview |
| `/monitoring` | Live metrics |
| `/history` | History graphs |
| `/activity` | Event timeline |
| `/alerts` | Active alerts |
| `/rules` | Rule management |
| `/services` | Service management |
| `/resources` | Resource health |
| `/guardian` | Guardian controls and status |
| `/observer` | Observer status and environment |
| `/recovery` | Recovery execution history |
| `/automation` | Automation executions |
| `/settings` | Settings |

---

# 9. Firmware

## Guardian

Dedicated hardware recovery controller.

Responsibilities:

- Power relay
- Reset relay
- Power LED input
- HTTP API

Endpoints: `/health`, `/status`, `/power`, `/reset`

## Observer

Environmental monitoring subsystem.

Responsibilities:

- Temperature sensor
- Humidity sensor
- Local display
- HTTP API

Endpoints: `/health`, `/status`, `/environment`

---

# 10. Coding Standards

- Small functions
- Small packages
- Clear naming
- No duplicated logic
- No premature optimization

---

# 11. Git Rules

- Every commit must compile.
- Main branch must remain stable.
- Conventional Commits are mandatory.

---

# 12. Documentation Rules

Every major architectural decision should be recorded as an ADR.

---

# 13. Testing Rules

Code is not complete until it can be verified.

Every sprint must define its verification procedure before implementation begins.

---

# 14. Long-Term Goal

Sentinel is designed to outlive its first implementation.

The architecture should remain understandable, maintainable and extensible for many years.