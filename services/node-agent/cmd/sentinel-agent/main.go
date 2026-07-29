package main

import (
	"fmt"
	"os"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/application"
)

const (
	ApplicationName    = "Sentinel Node Agent"
	ApplicationVersion = "0.1.0-dev"
)

func main() {
	fmt.Printf("%s %s\n", ApplicationName, ApplicationVersion)

	app, err := application.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Application failed: %v\n", err)
		os.Exit(1)
	}
}
