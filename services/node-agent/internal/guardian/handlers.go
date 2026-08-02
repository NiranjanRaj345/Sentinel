package guardian

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
		case "/guardian/status":
			h.status(w, r)
		case "/guardian/health":
			h.health(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/guardian/power":
			h.power(w, r)
		case "/guardian/reset":
			h.reset(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) status(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	status, err := h.service.Status(ctx)
	if err != nil {
		h.log.Error("guardian status error: %v", err)
		http.Error(w, "failed to load guardian status", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handler) power(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req PowerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.Power(ctx, req.Action); err != nil {
		h.log.Error("guardian power error: %v", err)
		http.Error(w, "failed to execute power action", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"result":"accepted"}`))
}

func (h *Handler) reset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req ResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.Reset(ctx, req.Action); err != nil {
		h.log.Error("guardian reset error: %v", err)
		http.Error(w, "failed to execute reset action", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"result":"accepted"}`))
}
