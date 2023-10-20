package translator

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	loadEnv() // Load .env

	r := gin.Default()

	// Define an AWS session and translator service
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		log.Fatal(err)
	}
	svc := translate.New(sess)

	r.POST("/translate", func(c *gin.Context) {
		// Get the input text and target language from the request
		var request struct {
			Text           string `json:"text" binding:"required"`
			TargetLanguage string `json:"target_language" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		// Create the translation input
		input := &translate.TextInput{
			SourceLanguageCode: aws.String("en"), //We only use English right?
			TargetLanguageCode: aws.String(request.TargetLanguage),
			Text:               aws.String(request.Text),
		}

		// Translate the text using the AWS Translate service
		result, err := svc.Text(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Translation failed"})
			return
		}

		// Extract the translated text
		translatedText := *result.TranslatedText

		//Mock accuracy
		accuracy := 0.80

		c.JSON(http.StatusOK, gin.H{
			"input_text":      request.Text,
			"translated_text": translatedText,
			"accuracy":        accuracy,
		})
	})

	// Run the server
	port := ":8080"
	fmt.Printf("Server is running on port %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
