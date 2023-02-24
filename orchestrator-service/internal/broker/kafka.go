package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/segmentio/kafka-go"
)

const (
	host  = "localhost:9092"
	topic = "processed-images-service"
)

type broker struct {
	conn *kafka.Conn
}

func New(c *kafka.Conn) *broker {
	return &broker{
		conn: c,
	}
}

type Message struct {
	Link string `json:"link"`
}

func (b *broker) Write(mes []byte) error {
	_, err := b.conn.Write(mes)
	return err
}

func (b *broker) Read() (arr []byte, err error) {
	reader := newKafkaReader(host, topic)
	defer func() {
		err = reader.Close()
	}()

	m, err := reader.ReadMessage(context.Background())
	if err != nil && err == io.EOF {
		return m.Value, err
	}
	var mes []Message
	err = json.Unmarshal(m.Value, &mes)
	if err != nil {
		fmt.Printf("Error unmarshaling message:%s", string(m.Value))
		return nil, err
	}
	fmt.Println(mes)
	return m.Value, err
}

func newKafkaReader(ports, topics string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{ports},
		Topic:    topics,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
}
