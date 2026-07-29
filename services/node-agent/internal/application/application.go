package application

import (
	"context"
	"fmt"

	"os"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
)

const ConfigPath = "config.yaml"

// Application represents the Sentinel Node Agent.
type Application struct {
	cfg    config.Config
	logger *logger.Logger
	server *server.Server
}

// New creates a new Application instance.
func New() (*Application, error) {
	cfg, err := config.Load(ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log, err := logger.NewFromString(cfg.Logging.Level, os.Stdout)
	if err != nil {
		return nil, fmt.Errorf("initialize logger: %w", err)
	}

	systemService := service.NewSystemService()
	metricsService := service.NewMetricsService()

	srv := server.New(
		cfg.Server.Address(),
		log,
		systemService,
		metricsService,
	)

	return &Application{
		cfg:    cfg,
		logger: log,
		server: srv,
	}, nil
}

// Run starts the application.
func (a *Application) Start() error {
	a.logger.Info("Starting Sentinel Node Agent")
	a.logger.Info("Configuration loaded")
	a.logger.Info("HTTP server listening on %s", a.cfg.Server.Address())

	return a.server.Start()
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.logger.Info("Shutting down Sentinel Node Agent")
	return a.server.Shutdown(ctx)
}
