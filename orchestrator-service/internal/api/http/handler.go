package http

import (
	"github.com/DarkJediDJ/image-service/orchestrator-service/internal/service"
	"github.com/aws/aws-sdk-go/aws/session"
	kafka "github.com/segmentio/kafka-go"
)

type Handler struct {
	service service.IService
}

func NewHandler(c *kafka.Conn, s *session.Session) *Handler {
	return &Handler{
		service: service.NewService(c, s),
	}
}
