package notification

import (
	"errors"
)

type MessageService interface {
	SendMessage(message string) error
	ReadMessage() (*Message, error)
}

type Message struct {
	Body string
}

type BrokerType int

const (
	SQS BrokerType = iota
	Kinesis
)

func NewMessageService(broker BrokerType, config interface{}) (MessageService, error) {
	switch broker {
	case SQS:
		c, ok := config.(SQSConfig)
		if !ok {
			return nil, errors.New("invalid config for SQS")
		}
		return NewSQSClient(c)
	case Kinesis:
		c, ok := config.(KinesisConfig)
		if !ok {
			return nil, errors.New("invalid config for Kinesis")
		}
		return NewKinesisClient(c)
	default:
		return nil, errors.New("unknown notification type")
	}
}
