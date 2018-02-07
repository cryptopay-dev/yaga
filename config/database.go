package config

// Database base configuration:
type Database struct {
	Address  string `yaml:"address" validate:"required"`
	Database string `yaml:"database" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}
