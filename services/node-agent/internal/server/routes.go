package server

import (
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/dashboard"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server/handlers"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/stream"
)

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", handlers.Metrics(s.metricsService))
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/history/latest", handlers.HistoryLatest(s.store))
	mux.HandleFunc("/history", handlers.HistoryRange(s.store))
	mux.HandleFunc("/stream", stream.StreamHandler(s.hub, s.log))
	mux.HandleFunc("/dashboard/overview", dashboard.OverviewHandler(s.dashboard, s.log))
	mux.HandleFunc("/dashboard/stream", dashboard.StreamHandler(s.dashboardHub, s.log))
	mux.Handle("/dashboard/history", dashboard.NewHistoryHandler(s.store, s.log))
	mux.HandleFunc("/dashboard/capabilities", handlers.Capabilities(s.nodeService))
}
