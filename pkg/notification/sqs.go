package notification

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConfig struct {
	Region              string
	QueueUrl            string
	MaxNumberOfMessages int32
}

type SQSClient struct {
	service *sqs.Client
	config  SQSConfig
}

func NewSQSClient(c SQSConfig) (MessageService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.Region),
	)

	if err != nil {
		return nil, err
	}

	return &SQSClient{
		service: sqs.NewFromConfig(cfg),
		config:  c,
	}, nil
}

func (c *SQSClient) SendMessage(message string) error {
	_, err := c.service.SendMessage(context.TODO(), &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(c.config.QueueUrl),
	})

	return err
}

func (c *SQSClient) ReadMessage(messages chan<- *Message, errors chan<- error) {
	output, err := c.service.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.config.QueueUrl),
		MaxNumberOfMessages: c.config.MaxNumberOfMessages,
		VisibilityTimeout:   30,
		WaitTimeSeconds:     20,
	})

	if err != nil {
		errors <- fmt.Errorf("error receiving messages: %w", err)
		return
	}

	for _, message := range output.Messages {
		messages <- &Message{
			*message.MessageId,
			*message.Body,
			*message.ReceiptHandle,
		}
	}
}

func (c *SQSClient) DeleteMessage(id string) error {
	_, err := c.service.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      &c.config.QueueUrl,
		ReceiptHandle: aws.String(id),
	})

	return err
}
