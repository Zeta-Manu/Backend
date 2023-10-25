package routes

import (
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "github.com/Zeta-Manu/Backend/docs"
	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
	"github.com/Zeta-Manu/Backend/internal/api/translator"
	"github.com/Zeta-Manu/Backend/internal/api/uploadtoS3"
)

func InitRoutes(router *gin.Engine, dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, translateService *translate.Translate) {
	videoController := controllers.NewVideoController(dbAdapter, s3Adapter)
	fileUploader := uploadtoS3.NewFileUploader(s3Adapter)
	Trans := translator.NewTranslator(translateService)

	docs.SwaggerInfo.BasePath = "/api"
	api := router.Group("/api")
	{
		api.POST("/postVideo", videoController.PostVideo)
		api.POST("/uploadtoS3", fileUploader.UploadFile)
		api.POST("/translate", Trans.TranslateText)
	}

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
}
