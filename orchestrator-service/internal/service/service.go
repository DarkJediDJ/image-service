package service

import (
	"encoding/json"

	"github.com/DarkJediDJ/image-service/orchestrator-service/internal/broker"
	"github.com/DarkJediDJ/image-service/orchestrator-service/internal/cloud"
	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/kevinburke/go.uuid"
	"github.com/segmentio/kafka-go"
)

type service struct {
	bucket Cloud
	broker Broker
}

func NewService(c *kafka.Conn, s *session.Session) *service {
	return &service{broker: broker.New(c), bucket: cloud.New(s)}
}

func (s *service) Process(img []byte) (arr []byte, err error) {
	file, err := s.bucket.Upload(img, uuid.NewV4().String())
	if err != nil {
		return nil, err
	}

	mes, err := json.Marshal(broker.Message{Link: file})
	if err != nil {
		return nil, err
	}

	if err = s.broker.Write(mes); err != nil {
		return nil, err
	}
	return s.broker.Read()
}
