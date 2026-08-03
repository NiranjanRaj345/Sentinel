package notification

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
	if h.service == nil {
		http.Error(w, "notification service not initialized", http.StatusInternalServerError)
		return
	}

	notifications, err := h.service.Recent(100)
	if err != nil {
		h.log.Error("list notifications error: %v", err)
		http.Error(w, "failed to load notifications", http.StatusInternalServerError)
		return
	}
	if notifications == nil {
		notifications = []Notification{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(RecentResponse{Notifications: notifications})
}
