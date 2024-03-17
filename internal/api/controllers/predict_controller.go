package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/config"
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
	s3_links, err := c.uploadVideoToS3(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while processing the video"})
	}
	// Insert a record into the database
	err = c.insertToS3Table(file.Filename, s3_links)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting record into database"})
		return
	}

	// Send the video to the ML API
	err = c.sendToML(file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while sending video to ML API"})
		return
	}
}

func (c *PredictController) uploadVideoToS3(file *multipart.FileHeader) (string, error) {
	// Open the file
	uploadedFile, err := file.Open()
	if err != nil {
		return "", err
	}
	defer uploadedFile.Close()

	// Read the file data
	fileData, err := io.ReadAll(uploadedFile)
	if err != nil {
		return "", err
	}

	// Upload the file to S3
	err = c.s3Adapter.PutObject(file.Filename, fileData)
	if err != nil {
		return "", err
	}

	// Construct the S3 URL
	s3_links := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.NewAppConfig().S3.BucketName, file.Filename)
	return s3_links, nil
}

func (c *PredictController) insertToS3Table(filename string, s3_links string) error {
	// SQL query to insert a new record into the database
	query := "INSERT INTO S3_Table (sub, s3_links) VALUES (?, ?) ON DUPLICATE KEY UPDATE s3_links = JSON_ARRAY_APPEND(COALESCE(s3_links, JSON_ARRAY()), '$', ?)"

	// Execute the query with the filename and status
	_, err := c.dbAdapter.Exec(query, filename, s3_links, s3_links)
	if err != nil {
		return err
	}

	return nil
}

func (c *PredictController) sendToML(filename string) error {
	fmt.Printf("Sending video %s to ML API\n", filename)
	return nil
}
