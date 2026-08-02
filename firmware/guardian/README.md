# Guardian Firmware

ESP32 firmware for the Sentinel Guardian hardware recovery controller.

## Hardware

- ESP32 DevKit
- Relay 1 → Power button
- Relay 2 → Reset button
- Input 1 → Power LED

## API

| Method | Path     | Purpose           |
| ------ | -------- | ----------------- |
| GET    | /status  | Device status     |
| POST   | /power   | Pulse power relay |
| POST   | /reset   | Pulse reset relay |
| GET    | /health  | Firmware health   |

### POST /power

```json
{ "action": "press" }
```

```json
{ "action": "release" }
```

### POST /reset

```json
{ "action": "press" }
```

```json
{ "action": "release" }
```

### GET /status

```json
{
  "status": "online",
  "firmware": "0.1.0",
  "uptime": 123,
  "powerButton": false,
  "resetButton": false,
  "powerLed": true,
  "lastSeen": "2026-08-02T16:30:00Z"
}
```

## Build

```bash
pio run
```

## Deploy

```bash
pio run --target upload
pio device monitor
```
