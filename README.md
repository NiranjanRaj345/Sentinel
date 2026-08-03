# 🛡️ Sentinel

> **Always Connected. Always Recoverable. Always Yours.**

Sentinel is a self-hosted personal infrastructure platform designed to monitor, manage, protect, and recover computers and devices from anywhere in the world.

---

## Architecture

```text
                Internet
                    │
            Fedora Laptop
                    │
               Tailscale VPN
                    │
              Home Router
                    │
      ┌─────────────┼─────────────┐
      │             │             │
      ▼             ▼             ▼
Desktop Agent   Guardian ESP32  Observer ESP32
      │             │             │
      └─────────────┼─────────────┘
                    │
             Mission Dashboard
```

---

## Vision

Sentinel begins with a single Windows desktop and evolves into a unified platform capable of managing:

- Windows desktops
- Linux servers
- NAS systems
- Raspberry Pi devices
- ESP32 hardware
- Future AI and home-lab infrastructure

---

## Core Components

### Node Agent

Runs on every managed machine.

Responsible for:

- System information
- Health reporting
- Software actions
- Event publishing
- Rule evaluation
- Automation execution
- Recovery orchestration

---

### Guardian

Hardware recovery controller.

Responsible for:

- Power control
- Reset control
- Hardware recovery

---

### Observer

Environmental monitoring subsystem.

Responsible for:

- Temperature sensing
- Humidity sensing
- Local status display

---

### Sentinel Dashboard

The primary interface used to monitor and control the infrastructure.

---

## Project Status

### ✅ Completed

- Remote Access Foundation
- Sentinel Node Agent
- Mission Dashboard
- Authentication
- Events Framework
- Rules Engine
- Automation Engine
- Recovery Engine
- Guardian Integration
- Observer Integration

### 🚧 In Progress

- Notification Providers
- Multi-node Support
- Mission Control

---

## License

This project is currently private.