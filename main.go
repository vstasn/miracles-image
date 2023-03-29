package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"miracles-image/m/v2/config"
	"miracles-image/m/v2/core"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/getsentry/sentry-go"
	"github.com/streadway/amqp"
)

func initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.AppConfig.SentryDSN,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)
}

func main() {
	log.Println("Starting image resize service...")

	if config.AppConfig.SentryDSN != "" {
		initSentry()
	}

	connectRabbitMQ, err := amqp.Dial(config.AppConfig.AMQPServerUrl)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	_, err = channelRabbitMQ.QueueDeclare(config.AppConfig.AMQPQueueName, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	messages, err := channelRabbitMQ.Consume(
		config.AppConfig.AMQPQueueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no local
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		log.Println(err)
	}

	s3Client := config.S3Client

	if s3Client == nil {
		panic("S3Client is blank")
	}

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages...")

	forever := make(chan bool)

	go func() {
		for message := range messages {
			log.Printf(" > Received message: %s\n", message.MessageId)

			addTask := &core.Task{}

			err := json.Unmarshal(message.Body, addTask)
			if err != nil {
				log.Println(err)
				continue
			}

			taskResult, err := addTask.ResizeImages()
			if err != nil {
				log.Println(err)
				continue
			}

			for _, item := range taskResult {
				_, err = config.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
					Bucket:      aws.String(config.AppConfig.S3.AWSStorageBucketName),
					Key:         aws.String(item.Key),
					Body:        item.Body,
					ContentType: &item.ContentType,
				})

				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	<-forever
}
