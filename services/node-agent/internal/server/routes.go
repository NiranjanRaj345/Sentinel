package server

import (
	"net/http"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/auth"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/dashboard"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/operations"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/rules"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server/handlers"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/stream"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/automation"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/guardian"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/observer"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/recovery"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
)

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/healthz/live", handlers.Live)
	mux.Handle("/healthz/ready", handlers.Ready(&s.ready))

	mux.Handle("/metrics", s.authRead(handlers.Metrics(s.metricsService)))
	mux.Handle("/history/latest", s.authRead(handlers.HistoryLatest(s.store)))
	mux.Handle("/history", s.authRead(handlers.HistoryRange(s.store)))
	mux.Handle("/stream", s.authRead(stream.StreamHandler(s.hub, s.log)))
	mux.Handle("/dashboard/overview", s.authRead(dashboard.OverviewHandler(s.dashboard, s.log)))
	mux.Handle("/dashboard/stream", s.authRead(dashboard.StreamHandler(s.dashboardHub, s.log)))
	mux.Handle("/dashboard/history", s.authRead(dashboard.NewHistoryHandler(s.store, s.log)))
	mux.Handle("/dashboard/capabilities", s.authRead(handlers.Capabilities(s.nodeService)))

	mux.Handle("/events/recent", s.authRead(events.EventsRecent(s.eventsService)))
	mux.Handle("/rules", s.authRead(rules.NewHandler(s.rulesService)))

	mux.Handle("/operations", s.authOperate(operations.NewHandler(s.operationsService, s.log)))
	mux.Handle("/services", s.authOperate(services.NewHandler(s.servicesService, s.log)))
	mux.Handle("/resources", s.authRead(resources.NewHandler(s.resourcesService, s.log)))
	mux.Handle("/automation/executions", s.authRead(automation.NewHandler(s.automationService, s.log)))
	mux.Handle("/guardian/status", s.authRead(guardian.NewHandler(s.guardianService, s.log)))
	mux.Handle("/guardian/power", s.authOperate(guardian.NewPowerHandler(s.guardianService, s.log)))
	mux.Handle("/guardian/reset", s.authOperate(guardian.NewResetHandler(s.guardianService, s.log)))
	mux.Handle("/observer/status", s.authRead(observer.NewHandler(s.observerService, s.log)))
	mux.Handle("/observer/environment", s.authRead(observer.NewHandler(s.observerService, s.log)))
	mux.Handle("/recovery/execute", s.authOperate(recovery.NewHandler(s.recoveryService, s.log)))
	mux.Handle("/recovery/recent", s.authRead(recovery.NewHandler(s.recoveryService, s.log)))
	mux.Handle("/notifications/recent", s.authRead(notification.NewHandler(s.notificationsService, s.log)))
	mux.Handle("/notifications/test", s.authOperate(notification.NewTestHandler(s.notificationsService, s.log)))
}

func (s *Server) authRead(next http.Handler) http.Handler {
	if s.authStore == nil || s.authStore.Empty() {
		return next
	}
	return auth.Authenticate(s.authStore)(auth.Authorize(auth.PermissionRead)(next))
}

func (s *Server) authOperate(next http.Handler) http.Handler {
	if s.authStore == nil || s.authStore.Empty() {
		return next
	}
	return auth.Authenticate(s.authStore)(auth.Authorize(auth.PermissionOperate)(next))
}
