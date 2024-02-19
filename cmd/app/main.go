package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/gin-gonic/gin"

	//"github.com/gin-contrib/cors"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/identityprovider"
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

	s3Adapter, err := s3.NewS3Adapter(appConfig.S3.BucketName, appConfig.S3.Region, creds)
	if err != nil {
		log.Fatalf("Failed to connect to S3: %v", err)
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(appConfig.S3.Region),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}
	transAdapter := translate.New(awsSession)

	idpAdapter, err := identityprovider.NewCognitoAdapter(appConfig.Cognito.Region, appConfig.Cognito.UserPoolID, appConfig.Cognito.ClientID)
	if err != nil {
		log.Fatalf("Failed to connect to Cognito: %v", err)
	}

	// Create a Gin router
	r := gin.Default()

	// Initialize routes
	routes.InitRoutes(r, db, *s3Adapter, transAdapter)
	routes.InitUserRoutes(r, idpAdapter)
	routes.InitPredictRoutes(r, db, *s3Adapter, *appConfig)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
