package metrics

import (
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	gnet "github.com/shirou/gopsutil/v4/net"
	"os"
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
	CPU     CPUInfo     `json:"cpu"`
	Memory  MemoryInfo  `json:"memory"`
	Disks   []DiskInfo  `json:"disks"`
	Network NetworkInfo `json:"network"`
}

func GetInfo() (Info, error) {
	cpuInfo, err := getCPUInfo()
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

	hostname, err := os.Hostname()
	if err != nil {
		return Info{}, err
	}

	interfaces, err := gnet.Interfaces()
	if err != nil {
		return Info{}, err
	}

	var networkInterfaces []NetworkInterface

	for _, iface := range interfaces {
		var addresses []string

		for _, addr := range iface.Addrs {
			addresses = append(addresses, addr.Addr)
		}

		networkInterfaces = append(networkInterfaces, NetworkInterface{
			Name:      iface.Name,
			MAC:       iface.HardwareAddr,
			Addresses: addresses,
		})
	}

	ioCounters, err := gnet.IOCounters(false)
	if err != nil {
		return Info{}, err
	}

	var networkIO NetworkIO

	if len(ioCounters) > 0 {
		io := ioCounters[0]

		networkIO = NetworkIO{
			BytesSent:     io.BytesSent,
			BytesReceived: io.BytesRecv,
			PacketsSent:   io.PacketsSent,
			PacketsRecv:   io.PacketsRecv,
		}
	}

	info := Info{
		CPU: cpuInfo,
		Memory: MemoryInfo{
			TotalBytes:   vm.Total,
			UsedBytes:    vm.Used,
			UsagePercent: vm.UsedPercent,
		},
		Disks: disks,
		Network: NetworkInfo{
			Hostname:   hostname,
			Interfaces: networkInterfaces,
			IO:         networkIO,
		},
	}

	return info, nil
}
