package metrics

import "github.com/shirou/gopsutil/v4/process"

// getProcessInfo collects information about running processes.
func getProcessInfo() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	processes := make([]ProcessInfo, 0, len(procs))

	for _, proc := range procs {
		name, err := proc.Name()
		if err != nil {
			// Skip processes we cannot inspect.
			continue
		}

		processes = append(processes, ProcessInfo{
			PID:  proc.Pid,
			Name: name,
		})
	}

	return processes, nil
}
