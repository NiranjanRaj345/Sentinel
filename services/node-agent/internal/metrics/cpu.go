package metrics

import "github.com/shirou/gopsutil/v4/cpu"

// getCPUInfo collects CPU usage information.
func getCPUInfo() (CPUInfo, error) {
	usage, err := cpu.Percent(0, false)
	if err != nil {
		return CPUInfo{}, err
	}

	return CPUInfo{
		UsagePercent: usage[0],
	}, nil
}
