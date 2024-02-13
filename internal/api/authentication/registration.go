package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// UserRegistration struct represents the user registration data
type UserRegistration struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserLogin struct represents the user login data
type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// Define routes
	r.POST("/register", RegisterHandler)
	r.POST("/login", LoginHandler)

	// Run the server
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating AWS session"})
		return
	}

	// Create Cognito service client
	cognitoClient := cognitoidentityprovider.New(sess)

	// Parse request body
	var userRegistration UserRegistration
	if err := c.BindJSON(&userRegistration); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Register user
	err = registerUser(cognitoClient, os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_CLIENT_ID"), userRegistration)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error registering user"})
		return
	}

	c.JSON(200, gin.H{"message": "User registered successfully"})
}

// LoginHandler handles user login
func LoginHandler(c *gin.Context) {
	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating AWS session"})
		return
	}

	// Create Cognito service client
	cognitoClient := cognitoidentityprovider.New(sess)

	// Parse request body
	var userLogin UserLogin
	if err := c.BindJSON(&userLogin); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Login user
	tokens, err := loginUser(cognitoClient, os.Getenv("COGNITO_USER_POOL_ID"), os.Getenv("COGNITO_CLIENT_ID"), userLogin)
	if err != nil {
		fmt.Println("Login error:", err)
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT tokens
	accessToken, err := generateJWTToken(*tokens.AccessToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating access token"})
		return
	}
	idToken, err := generateJWTToken(*tokens.IdToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating id token"})
		return
	}
	refreshToken, err := generateJWTToken(*tokens.RefreshToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error generating refresh token"})
		return
	}

	c.JSON(200, gin.H{"access_token": accessToken, "id_token": idToken, "refresh_token": refreshToken})
}

func registerUser(client *cognitoidentityprovider.CognitoIdentityProvider, userPoolID, clientID string, user UserRegistration) error {
	// Register user in Cognito User Pool
	_, err := client.SignUp(&cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(clientID),
		Username: aws.String(user.Email),
		Password: aws.String(user.Password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(user.Name),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(user.Email),
			},
		},
	})

	if err != nil {
		return err
	}

	return nil
}

// Generate JWT token
func generateJWTToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{
		"token": tokenString,
		"exp":   time.Now().Add(time.Hour * 1).Unix(), // Token expiry time
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func loginUser(client *cognitoidentityprovider.CognitoIdentityProvider, userPoolID, clientID string, user UserLogin) (*cognitoidentityprovider.AuthenticationResultType, error) {
	// Log in user using InitiateAuth
	result, err := client.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(user.Email),
			"PASSWORD": aws.String(user.Password),
		},
		ClientId: aws.String(clientID),
	})

	if err != nil {
		return nil, err
	}

	return result.AuthenticationResult, nil
}

/*curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"Thaksin\", \"email\": \"test@example.com\", \"password\": \"Password@69\"}" http://localhost:8080/register
  curl -X POST -H "Content-Type: application/json" -d "{\"email\": \"test@example.com\", \"password\": \"Password@69\"}" http://localhost:8080/login
*/
