package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

var AppConfig = getConfig()

type Config struct {
	SentryDSN     string `envconfig:"SENTRY_DSN"`
	AMQPServerUrl string `envconfig:"AMQP_SERVER_URL" required:"true"`
	AMQPQueueName string `envconfig:"AMQP_QUEUE_NAME" default:"image:resize"`
	S3            struct {
		AWSAccessKey         string `envconfig:"AWS_ACCESS_KEY_ID" required:"true"`
		AWSSecretKey         string `envconfig:"AWS_SECRET_ACCESS_KEY" required:"true"`
		AWSStorageBucketName string `envconfig:"AWS_STORAGE_BUCKET_NAME" required:"true"`
		AWSS3RegionName      string `envconfig:"AWS_S3_REGION_NAME" default:"eu-north-1" required:"true"`
	}
}

func getConfig() Config {
	var cnf Config
	err := envconfig.Process("", &cnf)
	if err != nil {
		log.Fatal(err)
	}
	return cnf
}
