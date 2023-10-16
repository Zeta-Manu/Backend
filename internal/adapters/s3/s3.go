package s3

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type S3Adapter struct {
	Session *session.Session
	Bucket  string
}

func NewS3Adapter(region, bucket string, creds *credentials.Credentials) (*S3Adapter, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}
	return &S3Adapter{
		Session: sess,
		Bucket:  bucket,
	}, nil
}

func (s *S3Adapter) GetObject(key string) ([]byte, error) {
	svc := s3.New(s.Session)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *S3Adapter) PutObject(key string, data []byte) error {
	svc := s3.New(s.Session)
	input := &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}

	_, err := svc.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}
