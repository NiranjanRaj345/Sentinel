package application

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/alert"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/dashboard"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/scheduler"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/storage/sqlite"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/stream"
)

const ConfigPath = "config.yaml"

type Application struct {
	cfg          config.Config
	logger       *logger.Logger
	server       *server.Server
	scheduler    *scheduler.Scheduler
	store        *sqlite.Store
	hub          *stream.Hub
	engine       *alert.Engine
	dashboard    *dashboard.Service
	dashboardHub *dashboard.Hub
}

func New() (*Application, error) {
	cfg, err := config.Load(ConfigPath)
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
	metricsScheduler := scheduler.New(interval, schedulerLog, store, hub, engine)
	metricsService := service.NewMetricsService(metricsScheduler)
	dashboardService := dashboard.NewService(metricsScheduler, engine, cfg, dashboardHub)
	metricsScheduler.SetPublishDashboard(dashboardService.PublishOverview)

	srv := server.New(
		cfg.Server.Address(),
		log.Component("server"),
		systemService,
		metricsService,
		store,
		hub,
		dashboardService,
		dashboardHub,
	)

	return &Application{
		cfg:          cfg,
		logger:       appLog,
		server:       srv,
		scheduler:    metricsScheduler,
		store:        store,
		hub:          hub,
		engine:       engine,
		dashboard:    dashboardService,
		dashboardHub: dashboardHub,
	}, nil
}

func (a *Application) Start() error {
	a.logger.Info("Starting Sentinel Node Agent")
	a.logger.Info("Configuration loaded")

	a.hub.Start()
	a.dashboardHub.Start()

	if err := a.scheduler.Start(); err != nil {
		return fmt.Errorf("start metrics scheduler: %w", err)
	}

	return a.server.Start()
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down Sentinel Node Agent")
	a.scheduler.Stop()
	a.hub.Stop()
	a.dashboardHub.Stop()

	if err := a.store.Close(); err != nil {
		a.logger.Error("failed to close storage: %v", err)
	}

	return a.server.Shutdown(ctx)
}
