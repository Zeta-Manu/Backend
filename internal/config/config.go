package config

import (
	"os"
)

// The application configuration
type AppConfig struct {
	DatabaseURL string
}

// initializes and returns the application configuration
func NewAppConfig() *AppConfig {
	return &AppConfig{
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}
