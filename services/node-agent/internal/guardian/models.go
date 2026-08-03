package guardian

import "time"

type GuardianStatus string

const (
	GuardianOnline  GuardianStatus = "online"
	GuardianOffline GuardianStatus = "offline"
)

type PowerAction string

const (
	PowerActionPress   PowerAction = "press"
	PowerActionRelease PowerAction = "release"
)

type ResetAction string

const (
	ResetActionPress   ResetAction = "press"
	ResetActionRelease ResetAction = "release"
)

type StatusResponse struct {
	Status         GuardianStatus `json:"status"`
	Firmware       string         `json:"firmware"`
	Uptime         int64          `json:"uptime"`
	PowerButton    bool           `json:"powerButton"`
	ResetButton    bool           `json:"resetButton"`
	PowerLed       bool           `json:"powerLed"`
	LastSeen       time.Time      `json:"lastSeen"`
}

type PowerRequest struct {
	Action PowerAction `json:"action"`
}

type ResetRequest struct {
	Action ResetAction `json:"action"`
}
