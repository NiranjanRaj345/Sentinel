# Deployment Guide

## Prerequisites

- Go 1.26+ (for building from source)
- SQLite (embedded, no external dependency)
- Windows 10+, Linux systemd, or Docker

## Build

```bash
cd services/node-agent
go build -o sentinel-agent ./cmd/sentinel-agent
go build -o sentinel-service ./cmd/sentinel-service
```

## Configuration

### Initialize config

```bash
./sentinel-agent init
```

This creates `config.yaml` with sensible defaults.

### Edit config

```bash
notepad config.yaml
```

Key settings:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

logging:
  level: "info"
  file: "logs/sentinel.log"

notifications:
  enabled: true
  providers:
    telegram:
      enabled: false
      bot_token: ""
      chat_id: ""
```

## Windows

### Install as Service

```powershell
# From an Administrator PowerShell prompt
.\deploy\windows\install-service.ps1 -BinaryPath "C:\sentinel\sentinel-service.exe"
```

### Uninstall

```powershell
.\deploy\windows\uninstall-service.ps1
```

### Verify

```powershell
Get-Service sentinel-agent
Get-Content C:\sentinel\logs\sentinel.log -Wait
```

## Linux

### Install systemd service

```bash
sudo cp deploy/systemd/sentinel-agent.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable sentinel-agent
sudo systemctl start sentinel-agent
```

### Verify

```bash
systemctl status sentinel-agent
journalctl -u sentinel-agent -f
```

### Health checks

```bash
curl http://localhost:8080/healthz/live
curl http://localhost:8080/healthz/ready
```

## Docker

```bash
docker compose -f deploy/docker/docker-compose.yml up -d
```

### Verify

```bash
docker compose -f deploy/docker/docker-compose.yml ps
docker logs sentinel
```

## Health Endpoints

| Endpoint         | Purpose                          |
| ---------------- | -------------------------------- |
| `/health`        | Detailed health with agent info  |
| `/healthz/live`  | Liveness probe (always 200 OK)   |
| `/healthz/ready` | Readiness probe (200 when ready) |

## Logs

- Console: stdout/stderr
- File: set `logging.file` in `config.yaml`

Log rotation is not built in. Use external tools:
- Linux: `logrotate`
- Windows: Task Scheduler + PowerShell

## Upgrades

### Stop, replace, start

```bash
systemctl stop sentinel-agent
cp sentinel-agent /usr/local/bin/sentinel-agent
systemctl start sentinel-agent
```

SQLite database is preserved between upgrades.

## Troubleshooting

### Service won't start

Check logs:
```bash
journalctl -u sentinel-agent -n 100
```

### Port already in use

Change `server.port` in `config.yaml`.

### Permissions

Ensure the service user has read/write access to:
- `config.yaml`
- `sentinel.db`
- `logs/` directory
