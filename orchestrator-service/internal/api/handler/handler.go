package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/kevinburke/go.uuid"
	kafka "github.com/segmentio/kafka-go"

	"github.com/DarkJediDJ/image-service/orchestrator-service/internal/broker"
	cloud "github.com/DarkJediDJ/image-service/orchestrator-service/internal/cloud"
)

type Handler struct {
	conn   *kafka.Conn
	bucket *cloud.Cloud
}

func Init(c *kafka.Conn, s *session.Session) *Handler {
	return &Handler{
		conn:   c,
		bucket: cloud.Init(s),
	}
}

func (h *Handler) Process(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")

	img, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	file, err := h.bucket.Upload(img, uuid.NewV4().String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusFailedDependency)
		return
	}

	mes, err := json.Marshal(broker.Message{Link: file})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.conn.Write([]byte(mes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	kafkaResp, err := broker.KafkaRead()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(kafkaResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(resp)

	w.WriteHeader(http.StatusCreated)
}
