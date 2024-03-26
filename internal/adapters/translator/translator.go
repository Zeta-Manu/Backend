package translator

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/translate"
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
func (tc *Translator) TranslateText(Text, TargetLanguage string) (translatedText string, accuracy float32, err error) {

	// Create the translation input
	input := &translate.TextInput{
		SourceLanguageCode: aws.String("en"), // Assuming source language is English
		TargetLanguageCode: aws.String(TargetLanguage),
		Text:               aws.String(Text),
	}

	// Translate the text using the AWS Translate service
	result, err := tc.TranslateService.Text(input)
	if err != nil {
		return "", 0, err
	}

	// Extract the translated text
	translatedText = *result.TranslatedText

	// Mock accuracy
	accuracy = 0.80

	return translatedText, accuracy, err
}
