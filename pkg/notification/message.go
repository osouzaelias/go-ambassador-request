package notification

import (
	"errors"
)

type MessageService interface {
	SendMessage(message string) error
	ReadMessage(messages chan<- *Message, errors chan<- error)
	DeleteMessage(id string) error
}

type Message struct {
	ID            string
	Body          string
	ReceiptHandle string
}

type BrokerType int

const (
	SQS BrokerType = iota
)

func NewMessageService(broker BrokerType, config interface{}) (MessageService, error) {
	switch broker {
	case SQS:
		c, ok := config.(SQSConfig)
		if !ok {
			return nil, errors.New("invalid config for SQS")
		}
		return NewSQSClient(c)
	default:
		return nil, errors.New("unknown notification type")
	}
}
