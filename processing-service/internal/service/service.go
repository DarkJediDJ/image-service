package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/DarkJediDJ/image-service/processing-service/internal/broker"
	"github.com/DarkJediDJ/image-service/processing-service/internal/cloud"
	img "github.com/DarkJediDJ/image-service/processing-service/internal/image"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/segmentio/kafka-go"
)

const (
	host  = "localhost:9092"
	topic = "default-images-service"
)

type service struct {
	bucket Cloud
	broker Broker
}

func NewService(c *kafka.Conn, s *session.Session) *service {
	return &service{broker: broker.New(c), bucket: cloud.New(s)}
}

func (s *service) Listen() (arr []broker.Message, err error) {
	incMes := make(chan broker.Message)

	for {
		go func(c chan broker.Message) {
			message, err := s.broker.Read()
			if err != nil {
				fmt.Printf("Error occured while reading from kafka: %e", err)
			}
			c <- message
		}(incMes)

		mes := <-incMes
		if (mes == broker.Message{}) {
			continue
		}

		filePath, err := s.bucket.Download(mes.Link)

		resiszedImages, err := img.Resize(filePath)
		if err != nil {
			fmt.Printf("Error occured while resizing images: %e", err)
		}

		var links []broker.Message

		for _, f := range resiszedImages {
			image, err := os.ReadFile(f)
			if err != nil {
				fmt.Printf("Error occured while reading file: %e", err)
			}

			name := strings.Split(f, cloud.FilePath)

			link, err := s.bucket.Upload(image, name[1])
			if err != nil {
				fmt.Printf("Error occured while uploading file to s3: %e", err)
			}

			links = append(links, broker.Message{Link: link})
		}

		img.RemoveContents(cloud.FilePath)

		resp, err := json.Marshal(links)
		if err != nil {
			fmt.Printf("Error occured while marshaling file: %e", err)
		}

		if err = s.broker.Write(resp); err != nil {
			return nil, err
		}
	}
}
