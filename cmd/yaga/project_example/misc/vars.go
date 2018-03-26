package misc

import (
	"time"

	"github.com/cryptopay-dev/yaga/config"
)

const (
	// Dev env
	Dev = "dev"
	// Name of application
	Name = "project"
	// Usage of application
	Usage = "Project service"
)

var (
	// Version of application by default is 'dev'
	Version = "dev"
	// BuildTime of application by default - current time
	BuildTime = time.Now().Format(time.RFC3339)
	// Debug from ENV:
	Debug = config.GetString("level") == Dev
)
