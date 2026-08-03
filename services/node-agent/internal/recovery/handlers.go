package recovery

import (
	"encoding/json"
	"fmt"
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
		case "/recovery/execute":
			h.execute(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	case http.MethodGet:
		switch r.URL.Path {
		case "/recovery/recent":
			h.recent(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) execute(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PolicyID string `json:"policyId"`
		Target   string `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	execution, err := h.service.Execute(r.Context(), req.PolicyID, req.Target)
	if err != nil {
		h.log.Error("recovery execution error: %v", err)
		http.Error(w, fmt.Sprintf("recovery execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(execution)
}

func (h *Handler) recent(w http.ResponseWriter, r *http.Request) {
	limit := 100
	executions, err := h.service.Recent(limit)
	if err != nil {
		h.log.Error("recovery recent error: %v", err)
		http.Error(w, fmt.Sprintf("failed to load recovery executions: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"executions": executions})
}
