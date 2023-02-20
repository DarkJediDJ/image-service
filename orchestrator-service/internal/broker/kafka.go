package broker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Link string `json:"link"`
}

func KafkaRead() (arr []Message, err error) {
	reader := newKafkaReader("localhost:9092", "processed-images-topic")
	defer reader.Close()
	var message Message

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		err = json.Unmarshal(m.Value, &message)
		if err != nil {
			return []Message{}, err
		}

		fmt.Println(message.Link)
		arr = append(arr, message)
	}

	err = reader.Close()
	return arr, err
}

func newKafkaReader(ports, topics string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{ports},
		Topic:    topics,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}
