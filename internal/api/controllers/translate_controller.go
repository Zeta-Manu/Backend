package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/domain/entity"
	valueobjects "github.com/Zeta-Manu/Backend/internal/domain/valueObjects"
)

type TranslateController struct {
	translateAdapter *translator.TranslateAdapter
}

func NewTranslateController(translateAdapter *translator.TranslateAdapter) *TranslateController {
	return &TranslateController{
		translateAdapter: translateAdapter,
	}
}

// TranslateController godoc
// @Summary Translate text
// @Description Translates the provided text into the target language
// @Accept json
// @Produce json
// @Param body body entity.TranslateJson true "Translation request"
// @Success 200 {object} entity.ResponseWrapper{data=valueobjects.TranslateControllerOutput} "Successful operation"
// @Failure 400 {object} entity.ErrorWrapper "Bad request"
// @Failure 500 {object} entity.ErrorWrapper "Internal server error"
// @Router /translate [post]
func (tc *TranslateController) TranslateText(c *gin.Context) {
	const SOURCELANGUAGE = "en"
	var req entity.TranslateJson
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := tc.translateAdapter.TranslateText(*req.Text, SOURCELANGUAGE, *req.TargetLanguage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": valueobjects.TranslateControllerOutput{
		OriginalText:   req.Text,
		TranslatedText: result.TranslateText,
	}})
}
