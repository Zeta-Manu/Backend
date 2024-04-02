package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Zeta-Manu/Backend/docs"
	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/adapters/sagemaker"
	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/api/routes"
	"github.com/Zeta-Manu/Backend/internal/config"
)

// @title Manu Swagger API
// @version 1.0
// @description server

// @host localhost:8081
// @BasePath /api

func main() {
	// Initialize the application configuration
	appConfig := config.NewAppConfig()

	creds := credentials.NewStaticCredentials(appConfig.IAM.Key, appConfig.IAM.Secret, "")

	db, err := database.InitializeDatabase(appConfig.Database)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	s3Adapter, err := s3.NewS3Adapter(appConfig.S3.Region, appConfig.S3.BucketName, creds)
	if err != nil {
		log.Fatalf("Failed to connect to S3: %v", err)
	}

	sagemakerAdapter, err := sagemaker.NewSageMakerAdapter(appConfig.S3.Region)
	if err != nil {
		log.Fatalf("Failed to connect to Sagemaker: %v", err)
	}

	translateAdapter, err := translator.NewTranslateAdapter(creds)
	if err != nil {
		log.Fatalf("Failed to connect to AWS Translate: %v", err)
	}

	// Create a Gin router
	r := gin.Default()

	// CROS-Middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	// Initialize routes
	routes.InitTranslateRoutes(r, *translateAdapter)
	routes.InitPredictRoutes(r, db, *s3Adapter, *sagemakerAdapter, *translateAdapter, *appConfig)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
