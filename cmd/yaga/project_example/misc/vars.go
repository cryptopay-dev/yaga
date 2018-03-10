package misc

import (
	"os"
	"time"
)

const (
	// Dev env
	Dev = "dev"
	// Name of application
	Name = "project"
	// Usage of application
	Usage = "Project service"
	// Config file path
	Config = "config.yml"
)

var (
	// Version of application by default is 'dev'
	Version = "dev"
	// BuildTime of application by default - current time
	BuildTime = time.Now().Format(time.RFC3339)
	// Debug from ENV:
	Debug = os.Getenv("LEVEL") == Dev
)
