package node

type Capability string

const (
	CapabilityMonitoring    Capability = "monitoring"
	CapabilityRemoteDesktop Capability = "remote_desktop"
	CapabilityVPN           Capability = "vpn"
	CapabilityGuardian      Capability = "guardian"
	CapabilityObserver      Capability = "observer"
)

type CapabilityStatus struct {
	Capability Capability `json:"capability"`
	Available  bool       `json:"available"`
	State      string     `json:"state"`
	Details    string     `json:"details,omitempty"`
}

type CapabilitiesResponse struct {
	Capabilities []CapabilityStatus `json:"capabilities"`
}
