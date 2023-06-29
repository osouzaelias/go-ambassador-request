package acceptor

import (
	"fmt"
	"go-ambassador-request/pkg/notification"
	"os"
	"strconv"
)

type SQSMessageSender struct {
	client notification.MessageService
}

func NewSQSMessageSender() (*SQSMessageSender, error) {
	maxNumberMessages, err := strconv.ParseInt(os.Getenv("MAX_NUMBER_MESSAGES"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to convert environment variable MAX_NUMBER_MESSAGES to int32: %w", err)
	}

	config := notification.SQSConfig{
		Region:              os.Getenv("REGION"),
		QueueUrl:            os.Getenv("SQS_QUEUE_URL"),
		MaxNumberOfMessages: int32(maxNumberMessages),
	}

	client, err := notification.NewMessageService(notification.SQS, config)
	if err != nil {
		return nil, err
	}

	return &SQSMessageSender{client: client}, nil
}

func (s *SQSMessageSender) SendMessage(msg string) error {
	return s.client.SendMessage(msg)
}

func (s *SQSMessageSender) ReadMessage(messages chan<- *notification.Message, errors chan<- error) {
	s.client.ReadMessage(messages, errors)
}

func (s *SQSMessageSender) DeleteMessage(id string) error {
	return s.client.DeleteMessage(id)
}
