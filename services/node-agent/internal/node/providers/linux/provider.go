package providers

import (
	"os/exec"
	"strings"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/node"
)

type linuxProvider struct {
	log *logger.Logger
}

func NewLinuxProvider(log *logger.Logger) node.Provider {
	return &linuxProvider{log: log}
}

func (p *linuxProvider) Name() string {
	return "linux"
}

func (p *linuxProvider) Capabilities() []node.CapabilityStatus {
	return []node.CapabilityStatus{
		p.monitoring(),
		p.remoteDesktop(),
		p.vpn(),
		p.guardian(),
		p.observer(),
	}
}

func (p *linuxProvider) monitoring() node.CapabilityStatus {
	running := p.isServiceRunning("sentinel-agent")
	state := "active"
	if !running {
		state = "inactive"
	}
	return node.CapabilityStatus{
		Capability: node.CapabilityMonitoring,
		Available:  true,
		State:      state,
		Details:    "Sentinel Agent monitoring subsystem",
	}
}

func (p *linuxProvider) remoteDesktop() node.CapabilityStatus {
	running := p.isServiceRunning("sunshine")
	state := "ready"
	if !running {
		state = "unavailable"
	}
	return node.CapabilityStatus{
		Capability: node.CapabilityRemoteDesktop,
		Available:  running,
		State:      state,
		Details:    "Sunshine remote desktop server",
	}
}

func (p *linuxProvider) vpn() node.CapabilityStatus {
	running := p.isServiceRunning("tailscaled")
	state := "disconnected"
	if running {
		state = "connected"
	}
	return node.CapabilityStatus{
		Capability: node.CapabilityVPN,
		Available:  running,
		State:      state,
		Details:    "Tailscale VPN connectivity",
	}
}

func (p *linuxProvider) guardian() node.CapabilityStatus {
	running := p.isServiceRunning("guardian")
	state := "missing"
	if running {
		state = "active"
	}
	return node.CapabilityStatus{
		Capability: node.CapabilityGuardian,
		Available:  running,
		State:      state,
		Details:    "Guardian automation agent",
	}
}

func (p *linuxProvider) observer() node.CapabilityStatus {
	running := p.isServiceRunning("observer")
	state := "missing"
	if running {
		state = "active"
	}
	return node.CapabilityStatus{
		Capability: node.CapabilityObserver,
		Available:  running,
		State:      state,
		Details:    "Observer passive monitoring agent",
	}
}

func (p *linuxProvider) isServiceRunning(name string) bool {
	cmd := exec.Command("systemctl", "is-active", "--quiet", name)
	cmd.Dir = "/"
	err := cmd.Run()
	if err != nil {
		return false
	}

	output, err := exec.Command("systemctl", "show", "--property=SubState", name).Output()
	if err != nil {
		return false
	}

	parts := strings.SplitN(string(output), "=", 2)
	if len(parts) != 2 {
		return false
	}

	state := strings.TrimSpace(parts[1])
	return state == "running"
}
