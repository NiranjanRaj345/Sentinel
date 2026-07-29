package main

import (
	"fmt"
	"log"

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

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
