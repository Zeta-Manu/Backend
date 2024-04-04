package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
	"github.com/Zeta-Manu/Backend/internal/config"
	manu_auth "github.com/Zeta-Manu/manu-auth/pkg/middleware"
)

func InitPredictRoutes(router *gin.Engine, dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, translator translator.TranslateAdapter, cfg config.AppConfig) {
	predictController := controllers.NewPredictController(dbAdapter, s3Adapter, translator)

	user := router.Group("/api", manu_auth.AuthenticationMiddleware(cfg.JWT.JWTPublicKey))
	{
		user.POST("/predict", predictController.Predict)
	}
}
