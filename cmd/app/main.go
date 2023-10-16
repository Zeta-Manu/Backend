package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/api/routes"
	"github.com/Zeta-Manu/Backend/internal/config"
)

func main() {
	// Initialize the application configuration
	appConfig := config.NewAppConfig()

	creds := credentials.NewStaticCredentials(appConfig.IAM.Key, appConfig.IAM.Secret, "")

	db, err := database.InitializeDatabase(appConfig.Database)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	s3Adapter, err := s3.NewS3Adapter(appConfig.S3.BucketName, appConfig.S3.Region)
	if err != nil {
		log.Fatalf("Failed to connect to S3: %v", err)
	}

	// Create a Gin router
	r := gin.Default()

	// Initialize routes
	routes.InitRoutes(r, db, *s3Adapter)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
