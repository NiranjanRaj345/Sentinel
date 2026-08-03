package observer

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
	case http.MethodGet:
		switch r.URL.Path {
		case "/observer/status":
			h.status(w, r)
		case "/observer/environment":
			h.environment(w, r)
		case "/observer/health":
			h.health(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status, err := h.service.Status(ctx)
	if err != nil {
		h.log.Error("observer status error: %v", err)
		http.Error(w, "failed to load observer status", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func (h *Handler) environment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	env, err := h.service.Environment(ctx)
	if err != nil {
		h.log.Error("observer environment error: %v", err)
		http.Error(w, "failed to load observer environment", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(env)
}
