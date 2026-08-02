package guardian

import (
	"encoding/json"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type PowerHandler struct {
	service *Service
	log     *logger.Logger
}

func NewPowerHandler(service *Service, log *logger.Logger) *PowerHandler {
	return &PowerHandler{service: service, log: log}
}

func (h *PowerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PowerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.Power(r.Context(), req.Action); err != nil {
		h.log.Error("guardian power error: %v", err)
		http.Error(w, "failed to execute power action", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"result":"accepted"}`))
}
