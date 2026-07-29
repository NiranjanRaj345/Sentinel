package application

import "fmt"

// Application represents the Sentinel Node Agent application.
type Application struct{}

// New creates a new Application instance.
func New() *Application {
	return &Application{}
}

// Run starts the application.
func (a *Application) Run() error {
	fmt.Println("Application initialized.")
	return nil
}