package nodes

import "time"

type Status string

const (
	StatusOnline   Status = "online"
	StatusOffline  Status = "offline"
	StatusUnknown  Status = "unknown"
)

type Node struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Hostname  string    `json:"hostname"`
	Address   string    `json:"address"`
	Version   string    `json:"version"`
	Platform  string    `json:"platform"`
	Status    Status    `json:"status"`
	LastSeen  time.Time `json:"lastSeen"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Address  string `json:"address"`
	Version  string `json:"version"`
	Platform string `json:"platform"`
}

type HeartbeatRequest struct {
	ID string `json:"id"`
}

type NodesResponse struct {
	Nodes []Node `json:"nodes"`
}

func NewNode(req RegisterRequest) Node {
	return Node{
		ID:        req.ID,
		Name:      req.Name,
		Hostname:  req.Hostname,
		Address:   req.Address,
		Version:   req.Version,
		Platform:  req.Platform,
		Status:    StatusOnline,
		LastSeen:  time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}
}
