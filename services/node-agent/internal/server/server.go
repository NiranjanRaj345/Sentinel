package server

import (
	"context"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server/middleware"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
)

type Server struct {
	httpServer     *http.Server
	systemService  *service.SystemService
	metricsService *service.MetricsService
}

func New(
	addr string,
	log *logger.Logger,
	systemService *service.SystemService,
	metricsService *service.MetricsService,
) *Server {

	server := &Server{
		systemService:  systemService,
		metricsService: metricsService,
	}

	mux := http.NewServeMux()

	server.registerRoutes(mux)

	handler := middleware.Chain(
		mux,
		middleware.Recovery(log),
		middleware.RequestID,
		middleware.Logging(log),
	)

	server.httpServer = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	return server
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
