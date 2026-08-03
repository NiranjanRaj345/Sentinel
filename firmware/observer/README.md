# Observer Firmware

ESP32 firmware for the Sentinel Observer environmental monitoring subsystem.

## Hardware

- ESP32 DevKit
- I2C Display
- Temperature/Humidity Sensor

## API

| Method | Path           | Purpose                 |
| ------ | -------------- | ----------------------- |
| GET    | `/health`      | Firmware health         |
| GET    | `/status`      | Device status           |
| GET    | `/environment` | Temperature and humidity |

### GET /status

```json
{
  "status": "online",
  "firmware": "0.1.0",
  "uptime": 123
}
```

### GET /environment

```json
{
  "temperature": 28.5,
  "humidity": 64.0
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
