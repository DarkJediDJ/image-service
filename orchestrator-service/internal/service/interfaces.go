package service

type Service interface {
	Process(img []byte) (arr []byte, err error)
}

type Cloud interface {
	Upload(img []byte, key string) (string, error)
}

type Broker interface {
	Write(mes []byte) error
	Read() (arr []byte, err error)
}
