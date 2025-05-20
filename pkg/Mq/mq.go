package mq

import "fmt"

type MqHandler interface {
	Receive() error
	Close() error
	Handle(func([]byte) error)
}

func NewMqHandlerSvc(mq interface{}) (MqHandler, error) {

	switch mt := mq.(type) {
	case *RabbitMq:
		return mt, nil
	default:
		return nil, fmt.Errorf("invalid mq")
	}
}
