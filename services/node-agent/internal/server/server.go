package server

import (
	"context"
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/auth"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/dashboard"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/node"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/operations"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/rules"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server/middleware"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/stream"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/automation"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/guardian"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/observer"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/recovery"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

type Server struct {
	httpServer           *http.Server
	log                  *logger.Logger
	systemService        *service.SystemService
	metricsService       *service.MetricsService
	store                *sqlite.Store
	hub                  *stream.Hub
	dashboard            *dashboard.Service
	dashboardHub         *dashboard.Hub
	nodeService          *node.Service
	operationsService    *operations.Service
	authStore            *auth.TokenStore
	eventsService        *events.Service
	rulesService         *rules.Service
	servicesService      *services.Service
	resourcesService     *resources.Service
	automationService    *automation.Service
	guardianService      *guardian.Service
	observerService      *observer.Service
	recoveryService      *recovery.Service
	notificationsService *notification.Service
}

func New(
	addr string,
	log *logger.Logger,
	systemService *service.SystemService,
	metricsService *service.MetricsService,
	store *sqlite.Store,
	hub *stream.Hub,
	dashboard *dashboard.Service,
	dashboardHub *dashboard.Hub,
	nodeService *node.Service,
	operationsService *operations.Service,
	authStore *auth.TokenStore,
	eventsService *events.Service,
	rulesService *rules.Service,
	servicesService *services.Service,
	resourcesService *resources.Service,
	automationService *automation.Service,
	guardianService *guardian.Service,
	observerService *observer.Service,
	recoveryService *recovery.Service,
	notificationsService *notification.Service,
) *Server {

	server := &Server{
		log:                  log,
		systemService:        systemService,
		metricsService:       metricsService,
		store:                store,
		hub:                  hub,
		dashboard:            dashboard,
		dashboardHub:         dashboardHub,
		nodeService:          nodeService,
		operationsService:    operationsService,
		authStore:            authStore,
		eventsService:        eventsService,
		rulesService:         rulesService,
		servicesService:      servicesService,
		resourcesService:     resourcesService,
		automationService:    automationService,
		guardianService:      guardianService,
		observerService:      observerService,
		recoveryService:      recoveryService,
		notificationsService: notificationsService,
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
