package metrics

import "github.com/shirou/gopsutil/v4/disk"

// getDiskInfo collects disk usage information.
func getDiskInfo() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// Skip inaccessible partitions.
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

	return disks, nil
}
