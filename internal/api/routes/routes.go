package routes

import (
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
)

func InitRoutes(router *gin.Engine, dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, translateService *translate.Translate) {
	videoController := controllers.NewVideoController(dbAdapter, s3Adapter)
	fileUploader := controllers.NewFileUploader(s3Adapter)
	trans := translator.NewTranslator(translateService)

	api := router.Group("/api")
	{
		api.POST("/postVideo", videoController.PostVideo)
		api.POST("/uploadtoS3", fileUploader.UploadFile)
		api.POST("/translate", trans.TranslateText) // WARNING: Controllers missing!
	}
}
