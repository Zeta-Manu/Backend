package main

import (
	"github.com/Zeta-Manu/Backend/internal/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize routes
	routes.InitRoutes(r)

	r.Run()
}
