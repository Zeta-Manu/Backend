package translator

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"

	valueobjects "github.com/Zeta-Manu/Backend/internal/domain/valueObjects"
)

type TranslateAdapter struct {
	Client *translate.Translate
}

func NewTranslateAdapter(region string, creds *credentials.Credentials) (*TranslateAdapter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}

	client := translate.New(sess)
	return &TranslateAdapter{
		Client: client,
	}, nil
}

func (ta *TranslateAdapter) TranslateText(text string, sourceLanguage string, targetLanguage string) (*valueobjects.TranslateOutput, error) {
	input := &translate.TextInput{
		SourceLanguageCode: aws.String(sourceLanguage),
		TargetLanguageCode: aws.String(targetLanguage),
		Text:               aws.String(text),
	}

	result, err := ta.Client.Text(input)
	if err != nil {
		return nil, err
	}
	return &valueobjects.TranslateOutput{
		TranslateText: result.TranslatedText,
		Meta: &valueobjects.TranslateMeta{
			SourceLanguage: result.SourceLanguageCode,
			TargetLanguage: result.TargetLanguageCode,
		},
	}, nil
}
