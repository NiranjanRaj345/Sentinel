package main

import "fmt"

const (
	ApplicationName    = "Sentinel Node Agent"
	ApplicationVersion = "0.1.0-dev"
)

func main() {
	fmt.Printf("%s %s\n", ApplicationName, ApplicationVersion)
	fmt.Println("Starting Sentinel Node Agent...")
}