package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Zeta-Manu/Backend/internal/adapters/s3"
)

type FileUploader struct {
	S3Adapter s3.S3Adapter
}

func NewFileUploader(s3Adapter s3.S3Adapter) *FileUploader {
	return &FileUploader{
		S3Adapter: s3Adapter,
	}
}

//	@Summary		Save video to S3
//	@Produce		json
//	@Description	Uploads a file to S3 bucket
//	@Accept			mpfd
//	@Produce		json
//	@Param			file	formData	file	true	"File to upload"
//	@Success		200		{string}	string	"File uploaded successfully"
//	@Failure		400		{object}	object	"Bad Request"
//	@Failure		500		{object}	object	"Internal Server Error"
//	@Router			/uploadtoS3 [post]
func (fc *FileUploader) UploadFile(c *gin.Context) {
	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Open the file
	uploadedFile, openErr := file.Open()
	if openErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": openErr.Error()})
		return
	}
	defer uploadedFile.Close()

	// Read the file data
	fileData, readErr := ioutil.ReadAll(uploadedFile)
	if readErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": readErr.Error()})
		return
	}

	// Upload the file to S3
	uploadErr := fc.S3Adapter.PutObject(file.Filename, fileData)
	if uploadErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": uploadErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}
