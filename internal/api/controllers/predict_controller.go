package controllers

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/gin-gonic/gin"
)

// NOTE: Mp4 -> S3, S3_TABLE -> ML API
type PredictController struct {
	dbAdapter database.DBAdapter
	s3Adapter s3.S3Adapter
}

func NewPredictController(dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter) *PredictController {
	return &PredictController{
		dbAdapter: dbAdapter,
		s3Adapter: s3Adapter,
	}
}

// @Summary Upload a video for prediction
// @Description Uploads a video file to S3 and prepares it for machine learning prediction
// @Tags api
// @Security BearerAuth
// @SecurityDefinition BearerAuth
// @SecurityDefinition.In header
// @SecurityDefinition.Name Authorization
// @SecurityDefinition.Type apiKey
// @Accept  multipart/form-data
// @Produce  json
// @Param   video formData file true "Video file to upload"
// @Success  200 {object} map[string]interface{}
// @Failure  400 {object} map[string]interface{}
// @Router /predict [post]
func (c *PredictController) Predict(ctx *gin.Context) {
	// Get the uploaded video file from the request
	file, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
		return
	}
	err = c.uploadVideoToS3(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while processing the video"})
	}
}

func (c *PredictController) uploadVideoToS3(file *multipart.FileHeader) error {
	// Open the file
	uploadedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	// Read the file data
	fileData, err := io.ReadAll(uploadedFile)
	if err != nil {
		return err
	}

	// Upload the file to S3
	err = c.s3Adapter.PutObject(file.Filename, fileData)
	if err != nil {
		return err
	}

	return nil
}

func (c *PredictController) insertToS3Table() {
}

func (c *PredictController) sendToML() {
}
