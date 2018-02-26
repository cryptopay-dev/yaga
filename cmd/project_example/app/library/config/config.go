package config

// Config structure
type Config struct {
	Bind string `yaml:"bind" validate:"required"`
}
