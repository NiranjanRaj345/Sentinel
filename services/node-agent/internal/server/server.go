package server

import (
	"fmt"
	"net/http"
)

// Server owns the HTTP server lifecycle.
type Server struct {
	httpServer *http.Server
}

// New creates a new HTTP server.
func New(addr string) *Server {
	mux := http.NewServeMux()

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		httpServer: httpServer,
	}
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	fmt.Printf("Starting HTTP server on %s\n", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}