package server

import (
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server/handlers"
)

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", handlers.Metrics)
	mux.HandleFunc("/health", handlers.Health)
}
