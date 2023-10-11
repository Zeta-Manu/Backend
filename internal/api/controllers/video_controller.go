package controllers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/database"
)

type VideoController struct {
	db database.DBAdapter
}

func NewVideoController(db database.DBAdapter) *VideoController {
	return &VideoController{
		db: db,
	}
}

//	@Summary	Upload a video
//	@Produce	json
//	@Param		video	formData	file	true	"Vidoe File"
//	@Router		/video [post]
func (vc *VideoController) PostVideo(c *gin.Context) {
	// Get the uploaded video file from the request
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No video file provided"})
		return
	}
	defer file.Close()

	// Create a temporary directory to store the videos
	tempDir := "uploads"
	if err := os.MkdirAll(tempDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a directory"})
		return
	}

	// Generate a unique filename for the video
	uniqueFilename := generateUniqueFilename(header.Filename)
	videoFilePath := filepath.Join(tempDir, uniqueFilename)

	// Create a new file to store the uploaded video
	out, err := os.Create(videoFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create a video file"})
		return
	}
	defer out.Close()

	// Copy the uploaded video to the new file
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save the video file"})
		return
	}

	// Simulate sending the video to the ML model and include the mock ML response
	mlResponse, err := forwardVideoToML(videoFilePath)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to send the video to the ML model"},
		)
		return
	}

	// Respond with the mock ML model's interpretation result
	c.JSON(http.StatusOK, mlResponse)
}

func generateUniqueFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := filename[:len(filename)-len(ext)]
	return fmt.Sprintf("%s_%s%s", base, randomString(5), ext)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func forwardVideoToML(videoFilePath string) (gin.H, error) {
	// Simulate ML processing and generate a mock response
	mockMLResponse := gin.H{
		"result": "interpretation result for " + filepath.Base(videoFilePath),
	}
	return mockMLResponse, nil
}
