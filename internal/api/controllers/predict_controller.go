package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/adapters/sagemaker"
	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/config"
	valueobjects "github.com/Zeta-Manu/Backend/internal/domain/valueObjects"
)

// NOTE: Mp4 -> S3, S3_TABLE -> ML API
type PredictController struct {
	sageMakerAdapter sagemaker.SageMakerAdapter
	dbAdapter        database.DBAdapter
	s3Adapter        s3.S3Adapter
	translator       *translator.Translator
}

func NewPredictController(dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, sagemakerAdapter sagemaker.SageMakerAdapter, translator *translator.Translator) *PredictController {
	return &PredictController{
		dbAdapter:        dbAdapter,
		s3Adapter:        s3Adapter,
		sageMakerAdapter: sagemakerAdapter,
		translator:       translator,
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
// @Param Authorization header string true "Bearer {token}"
// @Param   video formData file true "Video file to upload"
// @Success  200 {object} map[string]interface{}
// @Failure  400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /predict [post]
func (c *PredictController) Predict(ctx *gin.Context) {
	// Get the uploaded video file from the request
	file, err := ctx.FormFile("video")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
		return
	}
	s3Link, err := c.uploadVideoToS3(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while processing the video"})
	}
	// Insert a record into the database
	err = c.insertToS3Table(file.Filename, s3Link)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting record into database"})
		return
	}

	// Send the video to the ML API
	infer, err := c.sendToML(s3Link)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while sending video to ML API"})
		return
	}

	// Process the returned data from SageMaker
	processedData, err := c.processMLResult(infer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing ML result"})
		return
	}

	// Translate the processed data
	translatedData, err := c.translateData(processedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error translating data"})
		return
	}

	// Return the translated data to the client
	ctx.JSON(http.StatusOK, gin.H{"result": translatedData})
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
	s3Link := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.NewAppConfig().S3.BucketName, file.Filename)
	return s3Link, nil
}

func (c *PredictController) insertToS3Table(sub string, s3Link string) error {
	// SQL query to insert a new record into the database
	query := "INSERT INTO S3_Table (sub, s3_links) VALUES (?, ?) ON DUPLICATE KEY UPDATE s3_links = JSON_ARRAY_APPEND(COALESCE(s3_links, JSON_ARRAY()), '$', ?)"

	// Execute the query with the filename and status
	_, err := c.dbAdapter.Exec(query, sub, s3Link)
	if err != nil {
		return err
	}

	return nil
}

func (c *PredictController) sendToML(s3Link string) ([]byte, error) {
	input := valueobjects.SageMakerInput{
		Instance: []valueobjects.Instance{
			{
				Data: map[string]string{
					"s3Link": s3Link,
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(input)
	if err != nil {
		fmt.Println("Error marshaling input:", err)
		return nil, err
	}

	const (
		ENDPOINTNAME = "asl-deployment"
		CONTENTTYPE  = "application/json"
	)

	result, err := c.sageMakerAdapter.InvokeEndpoint(ENDPOINTNAME, CONTENTTYPE, jsonPayload)
	if err != nil {
		fmt.Println("Error invoking SageMaker endpoint:", err)
		return nil, err
	}
	return result, nil
}

func (c *PredictController) processMLResult(infer []byte) (string, error) {
	// The inference result is a JSON string
	var result map[string]interface{}
	err := json.Unmarshal(infer, &result)
	if err != nil {
		return "", err
	}

	// Process the result as needed
	label := result["label"].(string)

	return label, nil
}

func (c *PredictController) translateData(data string) (string, error) {
	// Translate the data using the translator instance
	translatedText, err := c.translator.TranslateText(data)
	if err != nil {
		return "", err
	}
	return translatedText, nil
}
