package utils

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

func UploadFileAws(key string, filename string) (*manager.UploadOutput, error) {
	path := "./uploads/" + filename

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("os.Open - filename: %s, err: %v", filename, err)
	}
	defer file.Close()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error : %v", err)
		return nil, err
	}
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)

	res, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("AWS_S3.S3_BUCKET")),
		Key:    aws.String(key),
		Body:   file,
		ACL:    "public-read",
	})

	return res, err
}
