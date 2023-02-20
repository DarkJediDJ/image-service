package service

import "github.com/DarkJediDJ/image-service/orchestrator-service/internal/broker"

type IService interface {
	Process(img []byte) (arr []broker.Message, err error)
}

type ICloud interface {
	Upload(img []byte, key string) (string, error)
}

type IBroker interface {
	Write(mes []byte) error
	Read() (arr []broker.Message, err error)
}
