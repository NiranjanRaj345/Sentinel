package server

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/health", s.handleHealth)
}