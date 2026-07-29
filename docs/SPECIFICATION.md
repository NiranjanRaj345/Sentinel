# 🛡️ Sentinel Specification

**Version:** 1.0.0-draft  
**Status:** Draft

---

# 1. Purpose

Sentinel is a self-hosted personal infrastructure platform designed to monitor, manage, protect, and recover computing devices from anywhere in the world.

The objective of Sentinel is to provide a reliable, secure, extensible, and maintainable platform that remains useful for many years.

---

# 2. Vision

Sentinel begins with a single desktop computer.

It will evolve into a unified platform capable of managing:

- Desktop computers
- Servers
- NAS systems
- Raspberry Pi devices
- ESP32 devices
- Home laboratory infrastructure
- Future AI infrastructure

---

# 3. Engineering Principles

Sentinel follows these principles.

## 3.1 Reliability First

Reliability is always preferred over additional features.

---

## 3.2 Simplicity

Prefer the simplest solution that correctly solves the problem.

---

## 3.3 Extensibility

Every component should be designed so future functionality can be added with minimal modification.

---

## 3.4 Single Responsibility

Every package, file and component should have exactly one responsibility.

---

## 3.5 Explicit Design

Architecture decisions must be documented before implementation.

---

## 3.6 Documentation

Important decisions must be recorded.

Documentation is considered part of the source code.

---

# 4. Repository Structure

```
Sentinel/

apps/
services/
firmware/
docs/
scripts/
```

Each top-level directory owns one domain of responsibility.

---

# 5. Core Components

## Dashboard

Primary user interface.

---

## Node Agent

Runs on managed machines.

Responsible for collecting system information and executing software actions.

---

## Guardian

Hardware recovery controller.

Responsible for power and reset operations.

---

## Gateway

Infrastructure coordinator.

Responsible for communication between Sentinel components.

---

# 6. Coding Standards

- Small functions
- Small packages
- Clear naming
- No duplicated logic
- No premature optimization

---

# 7. Git Rules

- Every commit must compile.
- Main branch must remain stable.
- Conventional Commits are mandatory.

---

# 8. Documentation Rules

Every major architectural decision should be recorded as an ADR.

---

# 9. Testing Rules

Code is not complete until it can be verified.

Every sprint must define its verification procedure before implementation begins.

---

# 10. Long-Term Goal

Sentinel is designed to outlive its first implementation.

The architecture should remain understandable, maintainable and extensible for many years.