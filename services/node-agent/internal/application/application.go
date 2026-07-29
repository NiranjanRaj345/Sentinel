package application

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/scheduler"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
)

const ConfigPath = "config.yaml"

type Application struct {
	cfg       config.Config
	logger    *logger.Logger
	server    *server.Server
	scheduler *scheduler.Scheduler
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

	systemService := service.NewSystemService()

	interval, err := time.ParseDuration(cfg.Metrics.Interval)
	if err != nil {
		return nil, fmt.Errorf("parse metrics interval: %w", err)
	}

	metricsScheduler := scheduler.New(interval, log.Component("scheduler"))
	metricsService := service.NewMetricsService(metricsScheduler)

	srv := server.New(
		cfg.Server.Address(),
		log.Component("server"),
		systemService,
		metricsService,
	)

	return &Application{
		cfg:       cfg,
		logger:    appLog,
		server:    srv,
		scheduler: metricsScheduler,
	}, nil
}

func (a *Application) Start() error {
	a.logger.Info("Starting Sentinel Node Agent")
	a.logger.Info("Configuration loaded")

	if err := a.scheduler.Start(); err != nil {
		return fmt.Errorf("start metrics scheduler: %w", err)
	}

	return a.server.Start()
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down Sentinel Node Agent")
	a.scheduler.Stop()
	return a.server.Shutdown(ctx)
}
