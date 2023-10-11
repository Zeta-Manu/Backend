package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/api/routes"
	"github.com/Zeta-Manu/Backend/internal/config"
)

func main() {
	// Initialize the application configuration
	appConfig := config.NewAppConfig()

	db, err := database.InitializeDatabase(appConfig.Database)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Create a Gin router
	r := gin.Default()

	// Initialize routes
	routes.InitRoutes(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
