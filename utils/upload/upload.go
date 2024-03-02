package upload

import (
	"fmt"
	"emailnotifl3n/app/config"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3UploaderInterface interface {
	UploadImage(fileHeader *multipart.FileHeader) (string, error)
	UploadMusic(fileHeader *multipart.FileHeader) (string, error)
}

type S3Uploader struct {
	sess *session.Session
}

func New() S3UploaderInterface {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(config.AWS_REGION),
		Credentials: credentials.NewStaticCredentials(
			config.AWS_ACCESS_KEY_ID,
			config.AWS_SECRET_ACCESS_KEY,
			""),
	})

	return &S3Uploader{sess: sess}
}

func (su *S3Uploader) UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	uploader := s3manager.NewUploader(su.sess)

	// Upload input configuration
	upParams := &s3manager.UploadInput{
		Bucket: aws.String("bucketl3n"),
		Key:    aws.String("foto/" + fileHeader.Filename),
		Body:   file,
	}

	// Run the upload
	resp, err := uploader.Upload(upParams)
	if err != nil {
		return "", fmt.Errorf("error uploading to S3: %w", err)
	}

	return resp.Location, nil
}

func (su *S3Uploader) UploadMusic(fileHeader *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".mp3" && ext != ".wav" && ext != ".flac" {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	uploader := s3manager.NewUploader(su.sess)

	// Upload input configuration
	upParams := &s3manager.UploadInput{
		Bucket: aws.String("bucketl3n"),
		Key:    aws.String("music/" + fileHeader.Filename),
		Body:   file,
	}

	// Run the upload
	resp, err := uploader.Upload(upParams)
	if err != nil {
		return "", fmt.Errorf("error uploading to S3: %w", err)
	}

	return resp.Location, nil
}
