package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/Zeta-Manu/Backend/internal/api/routes"
	"github.com/Zeta-Manu/Backend/internal/config"
	"github.com/Zeta-Manu/Backend/internal/database"
)

func main() {
	// Initialize the application configuration
	appConfig := config.NewAppConfig()

	dbDataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",appConfig.Database.User, appConfig.Database.Password, appConfig.Database.Host, appConfig.Database.Port, appConfig.Database.Name)
	db, err := database.NewDatabase(dbDataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Create a Gin router
	r := gin.Default()

	// Initialize routes
	routes.InitRoutes(r)

	r.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
