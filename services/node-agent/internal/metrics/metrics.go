package metrics

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
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

type DiskInfo struct {
	Device       string  `json:"device"`
	MountPoint   string  `json:"mountpoint"`
	FileSystem   string  `json:"filesystem"`
	TotalBytes   uint64  `json:"total_bytes"`
	UsedBytes    uint64  `json:"used_bytes"`
	FreeBytes    uint64  `json:"free_bytes"`
	UsagePercent float64 `json:"usage_percent"`
}

type Info struct {
	CPU    CPUInfo    `json:"cpu"`
	Memory MemoryInfo `json:"memory"`
	Disks  []DiskInfo `json:"disks"`
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

	partitions, err := disk.Partitions(false)
	if err != nil {
		return Info{}, err
	}

	var disks []DiskInfo

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		disks = append(disks, DiskInfo{
			Device:       partition.Device,
			MountPoint:   partition.Mountpoint,
			FileSystem:   partition.Fstype,
			TotalBytes:   usage.Total,
			UsedBytes:    usage.Used,
			FreeBytes:    usage.Free,
			UsagePercent: usage.UsedPercent,
		})
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
		Disks: disks,
	}

	return info, nil
}
