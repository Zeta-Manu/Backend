package routes

import (
	"TestBackend/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {
	videoController := controllers.NewVideoController()

	api := router.Group("/api")
	{
		api.POST("/postVideo", videoController.PostVideo)
	}
}
