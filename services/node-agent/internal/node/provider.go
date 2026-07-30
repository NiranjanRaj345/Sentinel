package node

type Provider interface {
	Name() string
	Capabilities() []CapabilityStatus
}
