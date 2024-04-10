package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type MLService interface {
	Predict(inputData string) ([]byte, error)
	CheckHealth() error
}

type mlServiceImpl struct {
	baseURL string
}

func NewMLService(baseURL string) (MLService, error) {
	mlService := &mlServiceImpl{
		baseURL: baseURL,
	}

	attempts := 3
	delay := 5 * time.Second

	// Perform a health check
	for i := 0; i < attempts; i++ {
		err := mlService.CheckHealth()
		if err == nil {
			return mlService, nil
		}
		log.Printf("Health check failed: %v", err)
		time.Sleep(delay)
		delay *= 2
	}

	return nil, fmt.Errorf("ML service is not healthy after %d attempts", attempts)
}

func (s *mlServiceImpl) CheckHealth() error {
	// Parse the endpoint URL
	u, err := url.Parse(s.baseURL)
	if err != nil {
		return err
	}

	// Add /healthz to the path
	u.Path += "/healthz"

	// Construct the full URL
	url := u.String()

	// Perform the HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ML service health check failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *mlServiceImpl) Predict(s3URI string) ([]byte, error) {
	// Parse the endpoint URL
	u, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, err
	}

	// Add /predict to the path
	u.Path += "/predict"

	// Add the S3 URI as a query parameter
	queryParams := u.Query()
	queryParams.Add("s3_uri", s3URI)
	u.RawQuery = queryParams.Encode()

	// Construct the full URL with the encoded query parameters
	url := u.String()

	// Perform the HTTP GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Return the response body as a byte slice
	return body, nil
}
