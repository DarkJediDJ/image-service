package service

import "github.com/DarkJediDJ/image-service/orchestrator-service/internal/broker"

type Service interface {
	Process(img []byte) (arr []broker.Message, err error)
}

type Cloud interface {
	Upload(img []byte, key string) (string, error)
}

type Broker interface {
	Write(mes []byte) error
	Read() (arr []broker.Message, err error)
}
