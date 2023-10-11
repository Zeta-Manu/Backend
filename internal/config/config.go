package config

import (
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// The application configuration
type AppConfig struct {
	Database DatabaseConfig
}

// initializes and returns the application configuration
func NewAppConfig() *AppConfig {
	dbConfig := DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}

	return &AppConfig{
		Database: dbConfig,
	}
}
