package server

import (
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
)

type Server struct {
	httpServer     *http.Server
	systemService  *service.SystemService
	metricsService *service.MetricsService
}

func New(
	addr string,
	systemService *service.SystemService,
	metricsService *service.MetricsService,
) *Server {

	server := &Server{
		systemService:  systemService,
		metricsService: metricsService,
	}

	mux := http.NewServeMux()

	server.registerRoutes(mux)

	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
