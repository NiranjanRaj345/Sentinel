package metrics

import "github.com/shirou/gopsutil/v4/mem"

// getMemoryInfo collects system memory information.
func getMemoryInfo() (MemoryInfo, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, err
	}

	return MemoryInfo{
		TotalBytes:   vm.Total,
		UsedBytes:    vm.Used,
		UsagePercent: vm.UsedPercent,
	}, nil
}
