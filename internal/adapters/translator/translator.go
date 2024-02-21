package translator

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/gin-gonic/gin"
)

type Translator struct {
	TranslateService *translate.Translate
}

func NewTranslator(translateService *translate.Translate) *Translator {
	return &Translator{
		TranslateService: translateService,
	}
}

// @Summary Translate text
// @Description Translates the provided text into the target language
// @Accept  json
// @Produce  json
// @Param text body string true "Text to translate"
// @Param target_language body string true "Target language code"
// @Success  200 {object} map[string]interface{} "{"input_text": "Input text", "translated_text": "Translated text", "accuracy":  0.80}"
// @Router /translate [post]
func (tc *Translator) TranslateText(c *gin.Context) {
	// Get the input text and target language from the request
	var request struct {
		Text           string `json:"text" binding:"required"`
		TargetLanguage string `json:"target_language" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Create the translation input
	input := &translate.TextInput{
		SourceLanguageCode: aws.String("en"), // Assuming source language is English
		TargetLanguageCode: aws.String(request.TargetLanguage),
		Text:               aws.String(request.Text),
	}

	// Translate the text using the AWS Translate service
	result, err := tc.TranslateService.Text(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Translation failed"})
		return
	}

	// Extract the translated text
	translatedText := *result.TranslatedText

	// Mock accuracy
	accuracy := 0.80

	c.JSON(http.StatusOK, gin.H{
		"input_text":      request.Text,
		"translated_text": translatedText,
		"accuracy":        accuracy,
	})
}
