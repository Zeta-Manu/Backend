package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/docs"
	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
	"github.com/Zeta-Manu/Backend/internal/config"
	manu_auth "github.com/Zeta-Manu/manu-auth/pkg/middleware"
)

func InitPredictRoutes(router *gin.Engine, dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, cfg config.AppConfig) {
	predictController := controllers.NewPredictController(dbAdapter, s3Adapter)

	// TODO: Fixed Swagger
	docs.SwaggerInfo.BasePath = "/api"
	user := router.Group("/api", manu_auth.AuthenticationMiddleware(cfg.JWT.JWTPublicKey))
	{
		user.POST("/predict", predictController.Predict)
	}
}
