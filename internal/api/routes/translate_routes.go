package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
)

func InitTranslateRoutes(router *gin.Engine, translateAdapter translator.TranslateAdapter) {
	translateController := controllers.NewTranslateController(&translateAdapter)

	translate := router.Group("/api")
	{
		translate.POST("/translate", translateController.TranslateText)
	}
}
