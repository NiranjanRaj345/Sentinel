package observer

import "time"

type ObserverStatus string

const (
	ObserverOnline  ObserverStatus = "online"
	ObserverOffline ObserverStatus = "offline"
)

type StatusResponse struct {
	Status    ObserverStatus `json:"status"`
	Firmware  string         `json:"firmware"`
	Uptime    int64          `json:"uptime"`
	LastSeen  time.Time      `json:"lastSeen"`
}

type EnvironmentResponse struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}
