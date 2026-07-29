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
	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/version"
)

func main() {
	fmt.Printf("%s %s\n",
		version.Build.Name,
		version.Build.Version,
	)

	app, err := application.New()
	if err != nil {
		log.Fatal(err)
	}

	// Start the application in the background.
	go func() {
		if err := app.Start(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Wait for an interrupt or termination signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("Shutdown signal received")

	// Allow up to 10 seconds for graceful shutdown.
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
