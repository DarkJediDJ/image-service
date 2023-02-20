package main

import (
	"context"
	"log"
	"os"

	server "github.com/DarkJediDJ/image-service/orchestrator-service/internal/api"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

const (
	addr = ":8080"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("An Error Occured while .env loading%v", err)
	}

	a := server.NewApp()

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "default-images-service", 0)
	if err != nil {
		log.Fatalf("An Error Occured while kafka connection%v", err)
	}

	region := os.Getenv("AWS_REGION")
	sess, err := session.NewSession(&aws.Config{
		Region:      &region,
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		log.Fatalf("An Error Occured while aws connection%v", err)
	}

	defer conn.Close()

	a.InitRouter(conn, sess)

	a.Run(addr)
}
