package routes

import (
	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
	"github.com/Zeta-Manu/Backend/internal/api/middleware"
	"github.com/Zeta-Manu/Backend/internal/config"
	"github.com/gin-gonic/gin"
)

func InitPredictRoutes(router *gin.Engine, dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, cfg config.AppConfig) {
	predictController := controllers.NewPredictController(dbAdapter, s3Adapter)

	// TODO: Fixed Swagger
	user := router.Group("/api", middleware.AuthenticationMiddleware(cfg))
	{
		user.POST("/predict", predictController.Predict)
	}
}
