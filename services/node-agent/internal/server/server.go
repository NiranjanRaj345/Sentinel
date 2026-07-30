package server

import (
	"context"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/dashboard"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server/middleware"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/stream"
)

type Server struct {
	httpServer     *http.Server
	log            *logger.Logger
	systemService  *service.SystemService
	metricsService *service.MetricsService
	store          *sqlite.Store
	hub            *stream.Hub
	dashboard      *dashboard.Service
}

func New(
	addr string,
	log *logger.Logger,
	systemService *service.SystemService,
	metricsService *service.MetricsService,
	store *sqlite.Store,
	hub *stream.Hub,
	dashboard *dashboard.Service,
) *Server {

	server := &Server{
		log:            log,
		systemService:  systemService,
		metricsService: metricsService,
		store:          store,
		hub:            hub,
		dashboard:      dashboard,
	}

	mux := http.NewServeMux()

	server.registerRoutes(mux)

	handler := middleware.Chain(
		mux,
		middleware.CORS,
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
	s.log.Info("HTTP server listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
