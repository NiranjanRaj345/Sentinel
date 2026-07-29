package metrics

import (
	"runtime"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/version"
)

type CPUInfo struct {
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryInfo struct {
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type DiskInfo struct {
	Device       string  `json:"device"`
	MountPoint   string  `json:"mountpoint"`
	FileSystem   string  `json:"filesystem"`
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type NetworkInterface struct {
	Name      string   `json:"name"`
	MAC       string   `json:"mac"`
	Addresses []string `json:"addresses"`
}

type NetworkIO struct {
	BytesSent     uint64 `json:"bytes_sent"`
	BytesReceived uint64 `json:"bytes_received"`
	PacketsSent   uint64 `json:"packets_sent"`
	PacketsRecv   uint64 `json:"packets_received"`
}

type NetworkInfo struct {
	Hostname   string             `json:"hostname"`
	Interfaces []NetworkInterface `json:"interfaces"`
	IO         NetworkIO          `json:"io"`
}

type Info struct {
	Metadata  Metadata      `json:"metadata"`
	CPU       CPUInfo       `json:"cpu"`
	Memory    MemoryInfo    `json:"memory"`
	Disks     []DiskInfo    `json:"disks"`
	Network   NetworkInfo   `json:"network"`
	Processes []ProcessInfo `json:"processes"`
}

func GetInfo() (Info, error) {
	start := time.Now()
	cpuInfo, err := getCPUInfo()
	if err != nil {
		return Info{}, err
	}
	memoryInfo, err := getMemoryInfo()
	if err != nil {
		return Info{}, err
	}
	diskInfo, err := getDiskInfo()
	if err != nil {
		return Info{}, err
	}
	networkInfo, err := getNetworkInfo()
	if err != nil {
		return Info{}, err
	}
	processInfo, err := getProcessInfo()
	if err != nil {
		return Info{}, err
	}
	end := time.Now()
	duration := end.Sub(start)

	return Info{
		Metadata: Metadata{
			Timestamp:            end,
			CollectionDurationMS: duration.Milliseconds(),
			Agent: AgentInfo{
				Name:         version.Build.Name,
				Version:      version.Build.Version,
				Platform:     runtime.GOOS,
				Architecture: runtime.GOARCH,
				GoVersion:    runtime.Version(),
			},
		},

		CPU:       cpuInfo,
		Memory:    memoryInfo,
		Disks:     diskInfo,
		Network:   networkInfo,
		Processes: processInfo,
	}, nil
}
