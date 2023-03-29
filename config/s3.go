package config

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client = getS3Client()

func getS3Client() *s3.Client {
	options := s3.Options{
		Region:      AppConfig.S3.AWSS3RegionName,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(AppConfig.S3.AWSAccessKey, AppConfig.S3.AWSSecretKey, "")),
	}

	client := s3.New(options)

	return client
}
