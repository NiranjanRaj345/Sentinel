package resources

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
		h.list(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	resources, err := h.service.List(ctx)
	if err != nil {
		h.log.Error("list resources error: %v", err)
		http.Error(w, "failed to list resources", http.StatusInternalServerError)
		return
	}
	if resources == nil {
		resources = []Resource{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resources)
}
