package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/YJU-OKURA/project_minori-gin-deployment-repo/constants"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Uploader interface {
	UploadImage(file *multipart.FileHeader, classID uint, isLogo bool) (string, error)
}

type awsUploader struct {
}

func NewAwsUploader() Uploader {
	return &awsUploader{}
}

// initializeS3Client S3クライアントを初期化
func initializeS3Client() (*s3.Client, error) {
	awsRegion := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_S3_ACCESS_KEY"),
			os.Getenv("AWS_S3_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", constants.DatabaseError, err)
	}
	return s3.NewFromConfig(cfg), nil
}

// UploadImage 画像をアップロード
func (u *awsUploader) UploadImage(fileHeader *multipart.FileHeader, classID uint, isLogo bool) (string, error) {
	log.Printf("UploadImage called with classID: %d, isLogo: %t", classID, isLogo)
	if fileHeader == nil {
		return "", fmt.Errorf(constants.ErrNoFileHeaderJP)
	}

	const maxSize = 10 << 20 // 10MB
	if fileHeader.Size > maxSize {
		return "", fmt.Errorf(constants.ErrFileSizeJP)
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		return "", fmt.Errorf("%s：%s", constants.ErrMimeTypeJP, mimeType)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrOpenFileJP, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			log.Printf("Failed to close file: %v", cerr)
			err = cerr
		}
	}()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%s: %w", constants.ErrReadFileDataJP, err)
	}

	s3Client, err := initializeS3Client()
	if err != nil {
		return "", err
	}
	uploader := manager.NewUploader(s3Client)

	// ファイル名の生成
	extension := filepath.Ext(fileHeader.Filename)
	uniqueFileName := fmt.Sprintf("images/%d/%s-%d%s", classID, strings.TrimSuffix(fileHeader.Filename, extension), time.Now().Unix(), extension)
	//log.Printf("Generated uniqueFileName: %s", uniqueFileName)

	// ロゴの場合、パスに 'logo/' を追加
	if isLogo {
		uniqueFileName = fmt.Sprintf("images/%d/logo/%s-%d%s", classID, strings.TrimSuffix(fileHeader.Filename, extension), time.Now().Unix(), extension)
	} else {
		uniqueFileName = fmt.Sprintf("images/%d/%s-%d%s", classID, strings.TrimSuffix(fileHeader.Filename, extension), time.Now().Unix(), extension)
	}

	bucketName := os.Getenv("AWS_S3_BUCKET_NAME")
	if bucketName == "" {
		return "", fmt.Errorf(constants.ErrLoadAWSConfigJP)
	}

	upParams := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(uniqueFileName),
		Body:   bytes.NewReader(fileData),
	}

	_, err = uploader.Upload(context.TODO(), upParams)
	if err != nil {
		log.Printf("Error in uploader.Upload: %v", err)
		return "", fmt.Errorf("%s: %w", constants.ErrUploadToS3JP, err)
	}

	cloudFrontURL := os.Getenv("AWS_CLOUDFRONT")
	if cloudFrontURL == "" {
		return "", fmt.Errorf(constants.ErrCloudFrontURLNotSetJP)
	}

	finalURL := fmt.Sprintf("%s/%s", cloudFrontURL, uniqueFileName)
	log.Printf("Final URL: %s", finalURL)
	return finalURL, nil
}
