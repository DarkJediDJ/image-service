package service

import (
	"github.com/DarkJediDJ/image-service/processing-service/internal/broker"
)

type Service interface {
	Resize(img []byte) (arr []broker.Message, err error)
}

type Cloud interface {
	Upload(img []byte, key string) (string, error)
	Download(key string) (string, error)
}

type Broker interface {
	Write(mes []byte) error
	Read() (mes broker.Message, err error)
}
