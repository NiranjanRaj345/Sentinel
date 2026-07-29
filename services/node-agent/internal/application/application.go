package application

import (
	"fmt"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/service"
)

const ConfigPath = "config.yaml"

// Application represents the Sentinel Node Agent.
type Application struct {
	cfg    config.Config
	server *server.Server
}

// New creates a new Application instance.
func New() (*Application, error) {
	cfg, err := config.Load(ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	systemService := service.NewSystemService()
	metricsService := service.NewMetricsService()

	server := server.New(
		cfg.Server.Address(),
		systemService,
		metricsService,
	)

	return &Application{
		cfg:    cfg,
		server: server,
	}, nil
}

// Run starts the application.
func (a *Application) Run() error {
	return a.server.Start()
}
