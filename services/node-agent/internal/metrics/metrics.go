package metrics

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type CPUInfo struct {
	UsagePercent float64 `json:"usage_percent"`
}

type MemoryInfo struct {
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type Info struct {
	CPU    CPUInfo    `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
}

func GetInfo() (Info, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return Info{}, err
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return Info{}, err
	}

	info := Info{
		CPU: CPUInfo{
			UsagePercent: cpuPercent[0],
		},
		Memory: MemoryInfo{
			TotalBytes:   vm.Total,
			UsedBytes:    vm.Used,
			UsagePercent: vm.UsedPercent,
		},
	}

	return info, nil
}
