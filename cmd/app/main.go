package main

import (
	"TestBackend/internal/api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Initialize routes
	routes.InitRoutes(r)

	r.Run()
}
