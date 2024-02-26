package utils

import (
	"bytes"
	"context"
	"fmt"
	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"mime/multipart"
	"os"
)

type Uploader interface {
	UploadImage(file *multipart.FileHeader) (string, error)
}

type awsUploader struct {
}

func NewAwsUploader() Uploader {
	return &awsUploader{}
}

// initializeS3Client S3クライアントを初期化
func initializeS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.DatabaseError, err) // Using constants for error message
	}
	return s3.NewFromConfig(cfg), nil
}

// UploadImage 画像をアップロード
func (u *awsUploader) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf(constants.ErrNoFileHeaderJP)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrOpenFileJP, err)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrReadFileDataJP, err)
	}

	s3Client, err := initializeS3Client()
	if err != nil {
		return "", err
	}
	uploader := manager.NewUploader(s3Client)

	upParams := &s3.PutObjectInput{
		Bucket: aws.String("images"), // Assuming bucket name is 'images'
		Key:    aws.String(fileHeader.Filename),
		Body:   bytes.NewReader(fileData),
	}

	_, err = uploader.Upload(context.TODO(), upParams)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrUploadToS3JP, err)
	}

	cloudFrontURL := os.Getenv("AWS_CLOUDFRONT")
	if cloudFrontURL == "" {
		return "", fmt.Errorf(constants.ErrCloudFrontURLNotSetJP)
	}

	finalURL := fmt.Sprintf("%s/%s", cloudFrontURL, fileHeader.Filename) // Ensure the URL is correctly formatted
	return finalURL, nil
}
