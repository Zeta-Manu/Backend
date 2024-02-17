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

type IAMConfig struct {
	Key    string
	Secret string
}

type S3Config struct {
	BucketName string
	Region     string
}

type CognitoConfig struct {
	UserPoolID string
	ClientID   string
	Region     string
}

// The application configuration
type AppConfig struct {
	Database DatabaseConfig
	IAM      IAMConfig
	S3       S3Config
	Cognito  CognitoConfig
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

	iamConfig := IAMConfig{
		Key:    os.Getenv("AWS_ACCESS_KEY_ID"),
		Secret: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	s3Config := S3Config{
		BucketName: os.Getenv("S3_BUCKET"),
		Region:     os.Getenv("REGION"),
	}

	cognitoConfig := CognitoConfig{
		UserPoolID: os.Getenv("COGNITO_POOL_ID"),
		ClientID:   os.Getenv("COGNITO_CLIENT_ID"),
		Region:     os.Getenv("REGION"),
	}

	return &AppConfig{
		Database: dbConfig,
		IAM:      iamConfig,
		S3:       s3Config,
		Cognito:  cognitoConfig,
	}
}
