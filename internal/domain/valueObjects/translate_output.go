package valueobjects

type TranslateOutput struct {
	TranslateText *string
	Meta          *TranslateMeta
}

type TranslateMeta struct {
	SourceLanguage *string
	TargetLanguage *string
}

type TranslateControllerOutput struct {
	OriginalText   *string
	TranslatedText *string
}
