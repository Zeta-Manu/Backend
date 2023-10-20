package uploadtoS3

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Initialize Gin
	r := gin.Default()

	//debug
	fmt.Println("REGION:", os.Getenv("REGION"))
	fmt.Println("S3_BUCKET:", os.Getenv("S3_BUCKET"))

	// Set up an AWS S3 session
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		panic(err)
	}
	uploader := s3manager.NewUploader(awsSession)

	// Define an endpoint to receive video files
	r.POST("/upload", func(c *gin.Context) {
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

		// Upload the file to S3
		_, uploadErr := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET")),
			Key:    aws.String(file.Filename),
			Body:   uploadedFile,
		})

		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": uploadErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	})

	// Run the Gin server
	r.Run(":8080")
}
