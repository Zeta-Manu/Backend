package routes

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/Zeta-Manu/Backend/docs"
	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
)

func InitRoutes(router *gin.Engine, dbAdapter database.DBAdapter) {
	videoController := controllers.NewVideoController(dbAdapter)

	docs.SwaggerInfo.BasePath = "/api"
	api := router.Group("/api")
	{
		api.POST("/postVideo", videoController.PostVideo)
	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
}
