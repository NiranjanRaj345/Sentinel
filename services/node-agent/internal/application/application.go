package application

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/alert"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/auth"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/dashboard"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events"
	eventSQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/events/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/node"
	nodeprovider "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/node/providers/linux"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/operations"
	opprovider "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/operations/providers/linux"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/rules"
	rulesSQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/rules/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/scheduler"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services"
	serviceslinux "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services/providers/linux"
	serviceswindows "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services/providers/windows"
	servicesSQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/services/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources"
	resourceslinux "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources/providers/linux"
	resourceswindows "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources/providers/windows"
	resourcesSQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/resources/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/stream"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/automation"
	automationSQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/automation/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/guardian"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/observer"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/recovery"
	recoverySQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/recovery/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification"
	notificationSQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification/storage/sqlite"
	telegramprovider "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/notification/providers/telegram"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/nodes"
	nodesSQLite "github.com/NiranjanRaj345/sentinel/services/node-agent/internal/nodes/storage/sqlite"
)

const ConfigPath = "config.yaml"

type Application struct {
	cfg                config.Config
	logger             *logger.Logger
	server             *server.Server
	scheduler          *scheduler.Scheduler
	store              *sqlite.Store
	hub                *stream.Hub
	engine             *alert.Engine
	dashboard          *dashboard.Service
	dashboardHub       *dashboard.Hub
	nodeService        *node.Service
	nodesService       *nodes.Service
	operationsService  *operations.Service
	eventsService      *events.Service
	rulesService       *rules.Service
	servicesService    *services.Service
	resourcesService   *resources.Service
	automationService  *automation.Service
	guardianService    *guardian.Service
	observerService    *observer.Service
	recoveryService    *recovery.Service
	notificationsService *notification.Service
	offlineCtx         context.Context
	offlineCancel      context.CancelFunc
}

func New() (*Application, error) {
	return NewWithConfig(ConfigPath)
}

func NewWithConfig(path string) (*Application, error) {
	cfg, err := config.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	log, err := logger.NewFromString(cfg.Logging.Level, os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	if cfg.Logging.File != "" {
		file, err := os.OpenFile(cfg.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("open log file: %w", err)
		}
		log = logger.New(logger.ParseLevelOrFallback(cfg.Logging.Level), file)
	}

	appLog := log.Component("application")
	schedulerLog := log.Component("scheduler")

	systemService := service.NewSystemService()

	interval, err := time.ParseDuration(cfg.Metrics.Interval)
	if err != nil {
		return nil, fmt.Errorf("parse metrics interval: %w", err)
	}

	store, err := sqlite.Open(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open storage: %w", err)
	}

	hub := stream.New(log)
	engine := alert.New(cfg.Alerts.Rules, schedulerLog)
	dashboardHub := dashboard.NewHub(log.Component("dashboard"))

	nodeProvider := nodeprovider.NewLinuxProvider(log.Component("node"))
	nodeService := node.NewService(nodeProvider, log.Component("node"))
	operationsProvider := opprovider.NewLinuxProvider(log.Component("operations"), nil)

	eventRepo, err := eventSQLite.OpenEvents(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open events storage: %w", err)
	}

	rulesRepo, err := rulesSQLite.OpenRules(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open rules storage: %w", err)
	}

	for _, rule := range rules.SeedRules() {
		if err := rulesRepo.Save(context.Background(), rule); err != nil {
			return nil, fmt.Errorf("seed rule %s: %w", rule.ID, err)
		}
	}

	servicesRepo, err := servicesSQLite.Open(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open services storage: %w", err)
	}

	servicesService := services.NewService(newServiceProvider(log.Component("services")), servicesRepo, log.Component("services"))

	eventsService := events.NewService(eventRepo, nil, log.Component("events"))

	nodesRepo, err := nodesSQLite.OpenNodes(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open nodes storage: %w", err)
	}

	nodesService := nodes.NewService(nodesRepo, eventsService.Publish, log.Component("nodes"))

	notificationsRepo, err := notificationSQLite.OpenNotifications(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open notifications storage: %w", err)
	}

	notificationsService := notification.NewService(notificationsRepo, eventsService.Publish, log.Component("notifications"))

	if cfg.Notifications.Enabled {
		for name, providerCfg := range cfg.Notifications.Providers {
			if !providerCfg.Enabled {
				continue
			}
			switch name {
			case "telegram":
				if providerCfg.BotToken != "" && providerCfg.ChatID != "" {
					telegramClient := telegramprovider.NewClient(providerCfg.BotToken, log.Component("telegram"))
					telegramProvider := telegramprovider.NewProvider(telegramClient, providerCfg.ChatID)
					notificationsService.AddProvider(telegramProvider)
				}
			}
		}
	}

	operationsService := operations.NewService(
		operationsProvider,
		operations.NewAuditor(log.Component("operations")),
		operations.NewValidator(operationsProvider),
		eventsService.Publish,
		notificationsService.Send,
		log.Component("operations"),
	)

	automationRepo, err := automationSQLite.Open(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open automation storage: %w", err)
	}

	guardianService := guardian.NewService(nil, eventsService.Publish, notificationsService.Send, log.Component("guardian"))

	automationEngine := automation.NewEngine(operationsService, guardianService, eventsService.Publish, notificationsService.Send, log.Component("automation"))
	automationService := automation.NewService(automationEngine, automationRepo, guardianService, notificationsService.Send, log.Component("automation"))

	recoveryRepo, err := recoverySQLite.OpenRecovery(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open recovery storage: %w", err)
	}

	recoveryExecutor := recovery.NewExecutor(guardianService, eventsService.Publish, log.Component("recovery"))
	recoveryEngine := recovery.NewEngine(recoveryExecutor, recoveryRepo, eventsService.Publish, notificationsService.Send, log.Component("recovery"))
	recoveryService := recovery.NewService(recoveryEngine, recoveryRepo, eventsService.Publish, notificationsService.Send, log.Component("recovery"))

	observerService := observer.NewService(nil, eventsService.Publish, notificationsService.Send, log.Component("observer"))

	rulesService := rules.NewService(rulesRepo, automationEngine, nil, log.Component("rules"))

	resourcesRepo, err := resourcesSQLite.Open(cfg.Storage.Path)
	if err != nil {
		return nil, fmt.Errorf("open resources storage: %w", err)
	}

	resourcesService := resources.NewService(newResourceProvider(log.Component("resources")), resourcesRepo, eventsService.Publish, log.Component("resources"))

	metricsScheduler := scheduler.New(
		interval,
		schedulerLog,
		store,
		hub,
		engine,
		eventsService.Publish,
		notificationsService.Send,
	)
	metricsService := service.NewMetricsService(metricsScheduler)
	dashboardService := dashboard.NewService(metricsScheduler, engine, cfg, dashboardHub)
	metricsScheduler.SetPublishDashboard(dashboardService.PublishOverview)

	authStore, err := auth.FromConfig(cfg.Auth)
	if err != nil {
		return nil, fmt.Errorf("initialize auth: %w", err)
	}

	srv := server.New(
		cfg.Server.Address(),
		log.Component("server"),
		systemService,
		metricsService,
		store,
		hub,
		dashboardService,
		dashboardHub,
		nodeService,
		nodesService,
		operationsService,
		authStore,
		eventsService,
		rulesService,
		servicesService,
		resourcesService,
		automationService,
		guardianService,
		observerService,
		recoveryService,
		notificationsService,
	)

	return &Application{
		cfg:                cfg,
		logger:             appLog,
		server:             srv,
		scheduler:          metricsScheduler,
		store:              store,
		hub:                hub,
		engine:             engine,
		dashboard:          dashboardService,
		dashboardHub:       dashboardHub,
		nodeService:        nodeService,
		nodesService:       nodesService,
		operationsService:  operationsService,
		eventsService:      eventsService,
		rulesService:       rulesService,
		servicesService:    servicesService,
		resourcesService:   resourcesService,
		automationService:  automationService,
		guardianService:    guardianService,
		observerService:    observerService,
		recoveryService:    recoveryService,
		notificationsService: notificationsService,
	}, nil
}

func newServiceProvider(log *logger.Logger) services.Provider {
	if runtime.GOOS == "windows" {
		return serviceswindows.NewWindowsProvider(log)
	}
	return serviceslinux.NewLinuxProvider(log)
}

func newResourceProvider(log *logger.Logger) resources.Provider {
	if runtime.GOOS == "windows" {
		return resourceswindows.NewWindowsProvider(log)
	}
	return resourceslinux.NewLinuxProvider(log)
}

func (a *Application) Start() error {
	a.logger.Info("Starting Sentinel Node Agent")
	a.logger.Info("Configuration loaded")
	a.logger.Info("Platform: %s", runtime.GOOS)

	a.hub.Start()
	a.dashboardHub.Start()

	if err := a.scheduler.Start(); err != nil {
		return fmt.Errorf("start metrics scheduler: %w", err)
	}

	timeout, err := time.ParseDuration(a.cfg.Nodes.HeartbeatTimeout)
	if err != nil {
		return fmt.Errorf("parse nodes heartbeat timeout: %w", err)
	}

	a.offlineCtx, a.offlineCancel = context.WithCancel(context.Background())
	go a.runOfflineChecker(a.offlineCtx, timeout)

	return a.server.Start()
}

func (a *Application) runOfflineChecker(ctx context.Context, timeout time.Duration) {
	ticker := time.NewTicker(timeout / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.nodesService.CheckOfflineNodes(ctx, timeout)
		}
	}
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down Sentinel Node Agent")
	a.scheduler.Stop()
	a.hub.Stop()
	a.dashboardHub.Stop()

	if a.offlineCancel != nil {
		a.offlineCancel()
	}

	if err := a.store.Close(); err != nil {
		a.logger.Error("failed to close storage: %v", err)
	}

	if err := a.eventsService.Close(); err != nil {
		a.logger.Error("failed to close events storage: %v", err)
	}

	if err := a.rulesService.Close(); err != nil {
		a.logger.Error("failed to close rules storage: %v", err)
	}

	if err := a.servicesService.Close(); err != nil {
		a.logger.Error("failed to close services storage: %v", err)
	}

	if err := a.resourcesService.Close(); err != nil {
		a.logger.Error("failed to close resources storage: %v", err)
	}

	if err := a.automationService.Close(); err != nil {
		a.logger.Error("failed to close automation storage: %v", err)
	}

	if err := a.recoveryService.Close(); err != nil {
		a.logger.Error("failed to close recovery storage: %v", err)
	}

	if err := a.notificationsService.Close(); err != nil {
		a.logger.Error("failed to close notifications storage: %v", err)
	}

	if err := a.nodesService.Close(); err != nil {
		a.logger.Error("failed to close nodes storage: %v", err)
	}

	return a.server.Shutdown(ctx)
}
