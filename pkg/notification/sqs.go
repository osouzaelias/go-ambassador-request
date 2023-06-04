package notification

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSConfig struct {
	Region   string
	QueueUrl string
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

func (c *SQSClient) ReadMessage() (*Message, error) {
	result, err := c.service.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.config.QueueUrl),
		MaxNumberOfMessages: 1,
		VisibilityTimeout:   30,
		WaitTimeSeconds:     20,
	})

	if err != nil {
		return nil, err
	}

	if len(result.Messages) == 0 {
		return nil, nil
	}

	msg := result.Messages[0]
	return &Message{
		Body: *msg.Body,
	}, nil
}
