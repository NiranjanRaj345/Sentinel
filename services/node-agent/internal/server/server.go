package server

import (
	"fmt"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
	"net/http"
)

// Server owns the HTTP server lifecycle.
type Server struct {
	mux           *http.ServeMux
	httpServer    *http.Server
	systemService *service.SystemService
}

// New creates a new HTTP server.
func New(
	addr string,
	systemService *service.SystemService,
) *Server {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	server := &Server{
		mux:           mux,
		httpServer:    httpServer,
		systemService: systemService,
	}

	server.registerRoutes()

	return server
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	fmt.Printf("Starting HTTP server on %s\n", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}
