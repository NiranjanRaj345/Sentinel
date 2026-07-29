package stream

import (
	"sync"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/metrics"
	"github.com/gorilla/websocket"
)

type Hub struct {
	log     *logger.Logger
	clients map[*websocket.Conn]struct{}

	latest    metrics.Info
	hasLatest bool

	mu   sync.RWMutex
	done chan struct{}
	once sync.Once
}

func New(log *logger.Logger) *Hub {
	return &Hub{
		log:     log,
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (h *Hub) Start() {
	h.done = make(chan struct{})
}

func (h *Hub) Stop() {
	h.once.Do(func() {
		close(h.done)
	})
}

func (h *Hub) Register(conn *websocket.Conn) {
	h.mu.Lock()
	h.clients[conn] = struct{}{}
	count := len(h.clients)
	h.mu.Unlock()

	h.log.Info("client connected")

	if h.hasLatest {
		if err := conn.WriteJSON(h.latest); err != nil {
			h.log.Error("failed to send initial snapshot to client: %v", err)
			_ = conn.Close()
			h.unregister(conn)
			return
		}
	}

	h.log.Info("broadcasted snapshot to %d client(s)", count)
}

func (h *Hub) Broadcast(info metrics.Info) {
	h.mu.Lock()
	h.latest = info
	h.hasLatest = true

	clients := make([]*websocket.Conn, 0, len(h.clients))
	for conn := range h.clients {
		clients = append(clients, conn)
	}
	h.mu.Unlock()

	for _, conn := range clients {
		if err := conn.WriteJSON(info); err != nil {
			h.log.Error("broadcast to client failed: %v", err)
			_ = conn.Close()
			h.unregister(conn)
		}
	}
}

func (h *Hub) unregister(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()

	h.log.Info("client disconnected")
}
