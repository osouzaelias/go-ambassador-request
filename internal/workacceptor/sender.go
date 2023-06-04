package workacceptor

import (
	"go-ambassador-request/pkg/notification"
	"os"
)

type SQSMessageSender struct {
	client notification.MessageService
}

func NewSQSMessageSender() (*SQSMessageSender, error) {
	config := notification.SQSConfig{
		Region:   os.Getenv("REGION"),
		QueueUrl: os.Getenv("SQS_QUEUE_URL"),
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
