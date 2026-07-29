package main

import (
	"fmt"
	"log"

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
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
