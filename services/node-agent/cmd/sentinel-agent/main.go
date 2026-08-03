package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/application"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/config"
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/version"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "init" {
		if err := runInit(); err != nil {
			log.Fatalf("init failed: %v", err)
		}
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "run" {
		runApplication(application.ConfigPath)
		return
	}

	configPath := application.ConfigPath
	if len(os.Args) > 1 && os.Args[1] == "--config" && len(os.Args) > 2 {
		configPath = os.Args[2]
	}

	fmt.Printf("%s %s\n",
		version.Build.Name,
		version.Build.Version,
	)

	runApplication(configPath)
}

func runApplication(configPath string) {
	app, err := application.NewWithConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := app.Start(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Sentinel Node Agent stopped gracefully")
}

func runInit() error {
	cfg := config.Default()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	path := application.ConfigPath
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("config file already exists: %s", path)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("Created %s\n", path)
	return nil
}
