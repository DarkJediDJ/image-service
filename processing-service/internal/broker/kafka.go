package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/segmentio/kafka-go"
)

type broker struct {
	conn   *kafka.Conn
	reader *kafka.Reader
}

func New(c *kafka.Conn) *broker {
	return &broker{
		conn:   c,
		reader: NewKafkaReader("localhost:9092", "default-images-service"),
	}
}

type Message struct {
	Link string `json:"link"`
}

func (b *broker) Write(mes []byte) error {
	_, err := b.conn.Write(mes)
	return err
}

func (b *broker) Read() (mes Message, err error) {
	m, err := b.reader.ReadMessage(context.Background())
	if err != nil && err == io.EOF {
		err = json.Unmarshal(m.Value, &mes)
		if err != nil {
			fmt.Printf("Error unmarshaling message:%s", string(m.Value))
			return Message{}, err
		}
		return mes, nil
	}

	err = json.Unmarshal(m.Value, &mes)
	if err != nil {
		fmt.Printf("Error unmarshaling message:%s", string(m.Value))
		return Message{}, err
	}

	return mes, err
}

func NewKafkaReader(ports, topics string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{ports},
		Topic:     topics,
		Partition: 0,
		MaxWait:   200 * time.Millisecond,
	})
}
