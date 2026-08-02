package guardian

import (
	"encoding/json"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type ResetHandler struct {
	service *Service
	log     *logger.Logger
}

func NewResetHandler(service *Service, log *logger.Logger) *ResetHandler {
	return &ResetHandler{service: service, log: log}
}

func (h *ResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.service.Reset(r.Context(), req.Action); err != nil {
		h.log.Error("guardian reset error: %v", err)
		http.Error(w, "failed to execute reset action", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(`{"result":"accepted"}`))
}
