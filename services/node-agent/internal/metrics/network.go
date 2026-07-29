package metrics

import (
	"os"

	gnet "github.com/shirou/gopsutil/v4/net"
)

// getNetworkInfo collects network interface and I/O statistics.
func getNetworkInfo() (NetworkInfo, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return NetworkInfo{}, err
	}

	interfaces, err := gnet.Interfaces()
	if err != nil {
		return NetworkInfo{}, err
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
		return NetworkInfo{}, err
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

	return NetworkInfo{
		Hostname:   hostname,
		Interfaces: networkInterfaces,
		IO:         networkIO,
	}, nil
}
