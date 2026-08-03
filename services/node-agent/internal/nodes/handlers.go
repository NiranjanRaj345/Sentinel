package nodes

import (
	"encoding/json"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Handler struct {
	service *Service
	log     *logger.Logger
}

func NewHandler(service *Service, log *logger.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/nodes/register":
			h.register(w, r)
		case "/nodes/heartbeat":
			h.heartbeat(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	case http.MethodGet:
		switch r.URL.Path {
		case "/nodes":
			h.list(w, r)
		default:
			if len(r.URL.Path) > 6 && r.URL.Path[:6] == "/nodes/" {
				h.detail(w, r)
				return
			}
			http.Error(w, "not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" || req.Name == "" {
		http.Error(w, "id and name are required", http.StatusBadRequest)
		return
	}

	node, err := h.service.Register(r.Context(), NewNode(req))
	if err != nil {
		h.log.Error("register node error: %v", err)
		http.Error(w, "failed to register node", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(node)
}

func (h *Handler) heartbeat(w http.ResponseWriter, r *http.Request) {
	var req HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	node, err := h.service.Heartbeat(r.Context(), req)
	if err != nil {
		h.log.Error("heartbeat error: %v", err)
		http.Error(w, "failed to process heartbeat", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(node)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.service.List(r.Context())
	if err != nil {
		h.log.Error("list nodes error: %v", err)
		http.Error(w, "failed to load nodes", http.StatusInternalServerError)
		return
	}
	if nodes == nil {
		nodes = []Node{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(NodesResponse{Nodes: nodes})
}

func (h *Handler) detail(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/nodes/"):]
	if id == "" {
		http.Error(w, "node id is required", http.StatusBadRequest)
		return
	}

	node, err := h.service.Get(r.Context(), id)
	if err != nil {
		h.log.Error("get node error: %v", err)
		http.Error(w, "failed to load node", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(node)
}
