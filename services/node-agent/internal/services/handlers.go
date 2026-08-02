package services

import (
	"encoding/json"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type Request struct {
	Action Action `json:"action"`
	Name   string `json:"name"`
}

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
	case http.MethodPost:
		h.execute(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	items, err := h.service.List(ctx)
	if err != nil {
		h.log.Error("list services error: %v", err)
		http.Error(w, "failed to list services", http.StatusInternalServerError)
		return
	}
	if items == nil {
		items = []ServiceItem{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Services []ServiceItem `json:"services"`
	}{Services: items})
}

func (h *Handler) execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Name == "" || (req.Action != ActionStart && req.Action != ActionStop && req.Action != ActionRestart) {
		http.Error(w, "invalid action or name", http.StatusBadRequest)
		return
	}
	result, err := h.service.Execute(ctx, req.Action, req.Name)
	if err != nil {
		h.log.Error("execute service error: %v", err)
		http.Error(w, "failed to execute service action", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
