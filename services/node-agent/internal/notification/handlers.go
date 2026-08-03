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

type TestRequest struct {
	Provider string `json:"provider,omitempty"`
}

type TestHandler struct {
	service *Service
	log     *logger.Logger
}

func NewTestHandler(service *Service, log *logger.Logger) *TestHandler {
	return &TestHandler{service: service, log: log}
}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.service == nil {
		http.Error(w, "notification service not initialized", http.StatusInternalServerError)
		return
	}

	var req TestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	notification := NewNotification(
		"test-"+r.RemoteAddr,
		"Test notification",
		"This is a test notification from Sentinel.",
		SeverityInfo,
		SourceSystem,
	)

	h.service.Send(r.Context(), notification)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "queued",
		"message": "test notification queued",
	})
}

