package sagemaker

import (
	"context"
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

func (a *SageMakerAdapter) InvokeEndpointAsync(ctx context.Context, endpointName, contentType string, payload []byte) (string, error) {
	input := &sagemakerruntime.InvokeEndpointAsyncInput{
		EndpointName: aws.String(endpointName),
		ContentType:  aws.String(contentType),
	}

	resp, err := a.client.InvokeEndpointAsyncWithContext(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke SageMaker endpoint asynchronously: %v", err)
	}

	return *resp.InferenceId, nil
}
