package server

import "net/http"

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/system", s.handleSystem)
	mux.HandleFunc("/metrics", s.handleMetrics)
}
