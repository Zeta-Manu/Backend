package controllers

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	httpadapter "github.com/Zeta-Manu/Backend/internal/adapters/http"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/domain/entity"
	valueobjects "github.com/Zeta-Manu/Backend/internal/domain/valueObjects"
)

// NOTE: Mp4 -> S3, S3_TABLE -> ML API
type PredictController struct {
	logger           *zap.Logger
	dbAdapter        database.DBAdapter
	s3Adapter        s3.S3Adapter
	translateAdapter translator.TranslateAdapter
	mlService        httpadapter.MLService
}

func NewPredictController(dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, translateAdapter translator.TranslateAdapter, mlService httpadapter.MLService, logger *zap.Logger) *PredictController {
	return &PredictController{
		dbAdapter:        dbAdapter,
		s3Adapter:        s3Adapter,
		translateAdapter: translateAdapter,
		logger:           logger,
		mlService:        mlService,
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
// @Param Authorization header string true "Bearer {token}" default(Bearer <Add access token here>)
// @Param   video formData file true "Video file to upload"
// @Success  200 {object} map[string]interface{}
// @Failure  400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /predict [post]
func (c *PredictController) Predict(ctx *gin.Context) {
	// Get the uploaded video file from the request
	file, err := ctx.FormFile("video")
	if err != nil {
		c.logger.Error("video form file failed", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
		return
	}

	s3Link, err := c.uploadVideoToS3(&c.s3Adapter.Bucket, file)
	if err != nil {
		c.logger.Error("Error uploading video to S3: ", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while processing the video"})
		return
	}

	sub, exists := ctx.Get("sub")
	if !exists {
		c.logger.Error("Cannot get subject", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Subject not found"})
		return
	}

	// Insert a record into the database
	err = c.insertToS3Table(sub.(string), s3Link)
	if err != nil {
		c.logger.Error("Error inserting record into database: ", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting record into database"})
		return
	}

	// Send the video to the ML API
	infer, err := c.sendToML(s3Link)
	if err != nil {
		c.logger.Error("Error sending video to ML API: ", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while sending video to ML API"})
		return
	}

	// Process the returned data from SageMaker
	avg, err := c.processMLResult(infer)
	if err != nil {
		c.logger.Error("Error processing ML result: ", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing ML result"})
		return
	}

	classes := getKeysFromProcessedAvgs(avg)
	responses := make([]entity.PredictResponse, len(classes))

	for i, class := range classes {
		// Translate the processed data
		translatedData, err := c.translateData(class, "TH")
		if err != nil {
			c.logger.Error("Error translating data: ", zap.Error(err))
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error translating data"})
			return
		}
		average := avg[i].Average
		sum := avg[i].Sum
		count := avg[i].Count

		responses[i] = entity.PredictResponse{
			Class:      class,
			Translated: *translatedData,
			Average:    average,
			Sum:        sum,
			Count:      count,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"result": responses})
}

func (c *PredictController) uploadVideoToS3(bucketName *string, file *multipart.FileHeader) (string, error) {
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
	s3Link := "s3://" + *bucketName + "/" + file.Filename
	// Debug: Log the constructed S3 URL
	return s3Link, nil
}

func (c *PredictController) insertToS3Table(sub string, s3Link string) error {
	// SQL query to insert a new record into the database
	query := "INSERT INTO S3_Table (sub, s3_links) VALUES (?, JSON_ARRAY(?))ON DUPLICATE KEY UPDATE s3_links = JSON_ARRAY_APPEND(COALESCE(s3_links, JSON_ARRAY()), '$', ?);"

	// Execute the query with the filename and status
	_, err := c.dbAdapter.Exec(query, sub, s3Link, s3Link)
	if err != nil {
		c.logger.Error("Failed to insert to db", zap.Error(err))
		return err
	}
	c.logger.Info("Insert to the table by sub", zap.String("sub", sub))

	return nil
}

func (c *PredictController) sendToML(s3Link string) ([]byte, error) {
	// Directly call the Predict method without using a goroutine
	result, err := c.mlService.Predict(s3Link)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *PredictController) processMLResult(infer []byte) ([]entity.ProcessedAvg, error) {
	var response valueobjects.MlResponse
	err := json.Unmarshal(infer, &response)
	if err != nil {
		c.logger.Error("Failed Unmarshal MlResponse", zap.Error(err))
		return nil, err
	}

	processedAvgs := make([]entity.ProcessedAvg, 0, len(response.Results.Avg))
	for key, value := range response.Results.Avg {
		processedAvgs = append(processedAvgs, entity.ProcessedAvg{
			Key:     key,
			Average: value.Average,
			Sum:     value.Sum,
			Count:   value.Count,
		})
	}
	return processedAvgs, nil
}

func (c *PredictController) translateData(data string, targetLanguage string) (*string, error) {
	const SOURCELANGUAGE = "en"
	translatedText, err := c.translateAdapter.TranslateText(data, SOURCELANGUAGE, targetLanguage)
	// Check for errors
	if err != nil {
		c.logger.Error("Cannot translate text input", zap.Error(err))
		return nil, err
	}

	return translatedText.TranslateText, nil
}

func getKeysFromProcessedAvgs(processedAvgs []entity.ProcessedAvg) []string {
	keys := make([]string, len(processedAvgs))
	for i, avg := range processedAvgs {
		keys[i] = avg.Key
	}
	return keys
}
