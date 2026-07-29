package application

import (
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/server"
)

// Application represents the Sentinel Node Agent.
type Application struct {
	server *server.Server
}

// New creates a new Application instance.
func New() *Application {
	return &Application{
		server: server.New(":8080"),
	}
}

// Run starts the application.
func (a *Application) Run() error {
	return a.server.Start()
}