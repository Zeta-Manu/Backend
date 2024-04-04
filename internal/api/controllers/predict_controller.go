package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
	"github.com/Zeta-Manu/Backend/internal/adapters/translator"
	"github.com/Zeta-Manu/Backend/internal/config"
	"github.com/Zeta-Manu/Backend/internal/domain/entity"
	valueobjects "github.com/Zeta-Manu/Backend/internal/domain/valueObjects"
)

// NOTE: Mp4 -> S3, S3_TABLE -> ML API
type PredictController struct {
	dbAdapter        database.DBAdapter
	s3Adapter        s3.S3Adapter
	translateAdapter translator.TranslateAdapter
}

func NewPredictController(dbAdapter database.DBAdapter, s3Adapter s3.S3Adapter, translateAdapter translator.TranslateAdapter) *PredictController {
	return &PredictController{
		dbAdapter:        dbAdapter,
		s3Adapter:        s3Adapter,
		translateAdapter: translateAdapter,
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

	appConfig := config.NewAppConfig()

	s3Link, err := c.uploadVideoToS3(&appConfig.S3.BucketName, &appConfig.S3.Region, file)
	if err != nil {
		fmt.Printf("Error uploading video to S3: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error while processing the video"})
	}

	// Debug: Log the S3 link of the uploaded video
	fmt.Printf("S3 link of uploaded video: %s\n", s3Link)

	// Insert a record into the database
	err = c.insertToS3Table(file.Filename, s3Link)
	if err != nil {
		fmt.Printf("Error inserting record into database: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting record into database"})
		return
	}

	// Send the video to the ML API
	infer, err := c.sendToML(appConfig.MLInference.ENDPOINT, s3Link)
	if err != nil {
		fmt.Printf("Error sending video to ML API: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error while sending video to ML API"})
		return
	}

	// Process the returned data from SageMaker
	avg, err := c.processMLResult(infer)
	if err != nil {
		fmt.Printf("Error processing ML result: %v\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing ML result"})
		return
	}

	classes := getKeysFromProcessedAvgs(avg)
	responses := make([]entity.PredictResponse, len(classes))

	for i, class := range classes {
		// Translate the processed data
		translatedData, err := c.translateData(class, "TH")
		if err != nil {
			fmt.Printf("Error translating data: %v\n", err)
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

	// Return the translated data to the client
	ctx.JSON(http.StatusOK, gin.H{"result": responses})
}

func (c *PredictController) uploadVideoToS3(bucketName *string, region *string, file *multipart.FileHeader) (string, error) {
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
	s3Link := "https://" + *bucketName + ".s3." + *region + ".amazonaws.com/" + file.Filename
	// Debug: Log the constructed S3 URL
	fmt.Printf("Constructed S3 URL: %s\n", s3Link)
	return s3Link, nil
}

func (c *PredictController) insertToS3Table(sub string, s3Link string) error {
	// SQL query to insert a new record into the database
	query := "INSERT INTO S3_Table (sub, s3_links) VALUES (?, JSON_ARRAY(?))ON DUPLICATE KEY UPDATE s3_links = JSON_ARRAY_APPEND(COALESCE(s3_links, JSON_ARRAY()), '$', ?);"

	// Execute the query with the filename and status
	_, err := c.dbAdapter.Exec(query, sub, s3Link, s3Link)
	if err != nil {
		return err
	}

	return nil
}

func (c *PredictController) sendToML(endpoint string, s3Link string) ([]byte, error) {
	url := endpoint + fmt.Sprintf("/predict?s3_uri=%s", s3Link)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error making the GET request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading the response body: %v\n", err)
		return nil, err
	}

	return body, nil
}

func (c *PredictController) processMLResult(infer []byte) ([]entity.ProcessedAvg, error) {
	var response valueobjects.MlResponse
	err := json.Unmarshal(infer, &response)
	if err != nil {
		log.Fatalf("Error Unmarshalling JSON: %v", err)
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
