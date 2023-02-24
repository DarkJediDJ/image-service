package main

import (
	"context"
	"log"
	"os"

	"github.com/DarkJediDJ/image-service/processing-service/internal/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

const (
	port  = ":3001"
	addr  = "localhost:9092"
	topic = "processed-images-service"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("An Error Occured while .env loading%v", err)
	}

	conn, err := kafka.DialLeader(context.Background(), "tcp", addr, topic, 0)
	if err != nil {
		log.Fatalf("An Error Occured while kafka connection%v", err)
	}

	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Credentials:      credentials.NewEnvCredentials(),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("http://localhost:4566"),
	})
	if err != nil {
		log.Fatalf("An Error Occured while aws connection%v", err)
	}

	defer conn.Close()

	s := service.NewService(conn, sess)
	s.Listen()
}
