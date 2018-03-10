package config

//import "github.com/cryptopay-dev/yaga/config"

// Config structure
type Config struct {
	Bind string `yaml:"bind" validate:"required"`
	// Database    config.Database `yaml:"database" validate:"required,dive"`
	// Redis    config.Redis `yaml:"redis" validate:"required,dive"`
}
