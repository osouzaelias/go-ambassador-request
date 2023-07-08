package acceptor

import (
	"encoding/json"
	"fmt"
	"go-ambassador-request/pkg/config"
	"go-ambassador-request/pkg/notification"
	"strconv"
)

type SQSMessageSender struct {
	client notification.MessageService
}

func NewSQSMessageSender(cfg *config.Config) (*SQSMessageSender, error) {
	maxNumberMessages, err := strconv.ParseInt(cfg.MaxNumberMessages, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert environment variable MAX_NUMBER_MESSAGES to int32: %w", err)
	}

	sqsConfig := notification.SQSConfig{
		Region:              cfg.Region,
		QueueUrl:            cfg.QueueUrl,
		MaxNumberOfMessages: int32(maxNumberMessages),
	}

	client, err := notification.NewMessageService(notification.SQS, sqsConfig)
	if err != nil {
		return nil, err
	}

	return &SQSMessageSender{client: client}, nil
}

func (s *SQSMessageSender) SendMessage(req Request) error {
	msg, _ := json.Marshal(req.Data)
	return s.client.SendMessage(string(msg))
}

func (s *SQSMessageSender) ReadMessage(message chan<- *notification.Message, err chan<- error) {
	s.client.ReadMessage(message, err)
}

func (s *SQSMessageSender) DeleteMessage(id string) error {
	return s.client.DeleteMessage(id)
}
