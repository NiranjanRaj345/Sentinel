package dashboard

import (
	"sync"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      Overview    `json:"data"`
}

type Hub struct {
	log     *logger.Logger
	clients map[*websocket.Conn]struct{}

	latest    Message
	hasLatest bool

	mu     sync.RWMutex
	done   chan struct{}
	once   sync.Once
}

func NewHub(log *logger.Logger) *Hub {
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
	if h.hasLatest {
		_ = conn.WriteJSON(h.latest)
	}
	h.mu.Unlock()

	h.log.Info("dashboard client connected")
}

func (h *Hub) Broadcast(message Message) {
	h.mu.Lock()
	h.latest = message
	h.hasLatest = true

	clients := make([]*websocket.Conn, 0, len(h.clients))
	for conn := range h.clients {
		clients = append(clients, conn)
	}
	h.mu.Unlock()

	for _, conn := range clients {
		if err := conn.WriteJSON(message); err != nil {
			h.log.Error("dashboard broadcast failed: %v", err)
			_ = conn.Close()
			h.unregister(conn)
		}
	}
}

func (h *Hub) unregister(conn *websocket.Conn) {
	h.mu.Lock()
	delete(h.clients, conn)
	h.mu.Unlock()

	h.log.Info("dashboard client disconnected")
}
