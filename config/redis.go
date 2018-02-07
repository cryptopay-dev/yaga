package config

// Redis default configuration
type Redis struct {
	Address  string `yaml:"address" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Database int    `yaml:"database" validate:"required"`
}
