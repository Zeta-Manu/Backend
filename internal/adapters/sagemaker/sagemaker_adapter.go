package sagemaker

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

type SageMakerAdapter struct {
	client *sagemakerruntime.SageMakerRuntime
}

func NewSageMakerAdapter(region string) (*SageMakerAdapter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}

	client := sagemakerruntime.New(sess)
	return &SageMakerAdapter{client: client}, nil
}

func (a *SageMakerAdapter) InvokeEndpoint(endpointName, contentType string, payload []byte) ([]byte, error) {
	input := &sagemakerruntime.InvokeEndpointInput{
		Body:         payload,
		ContentType:  aws.String(contentType),
		EndpointName: aws.String(endpointName),
	}

	result, err := a.client.InvokeEndpoint(input)
	if err != nil {
		return nil, fmt.Errorf("failed to invoke SageMaker endpoint: %v", err)
	}

	return result.Body, nil
}
