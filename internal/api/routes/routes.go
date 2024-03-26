package routes

import (
	"net/http"

	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/api/controllers"
)

func translateHandler(trans *translator.Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract data from the request, such as text and target language
		var requestData struct {
			Text           string `json:"text" binding:"required"`
			TargetLanguage string `json:"target_language" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Translate the text using the translator service
		translatedText, accuracy, err := trans.TranslateText(requestData.Text, requestData.TargetLanguage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Return the translated text and accuracy
		c.JSON(http.StatusOK, gin.H{
			"input_text":      requestData.Text,
			"translated_text": translatedText,
			"accuracy":        accuracy,
		})
	}
}
func InitRoutes(router *gin.Engine, dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, translateService *translate.Translate) {
	videoController := controllers.NewVideoController(dbAdapter, s3Adapter)
	fileUploader := controllers.NewFileUploader(s3Adapter)
	trans := translator.NewTranslator(translateService)

	api := router.Group("/api")
	{
		api.POST("/postVideo", videoController.PostVideo)
		api.POST("/uploadtoS3", fileUploader.UploadFile)
		api.POST("/translate", translateHandler(trans)) // WARNING: Controllers missing!
	}
}
