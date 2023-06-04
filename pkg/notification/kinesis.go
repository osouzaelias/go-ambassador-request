package notification

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

type KinesisConfig struct {
	Region     string
	StreamName string
	ShardId    string
}

type KinesisClient struct {
	service *kinesis.Client
	config  KinesisConfig
}

func NewKinesisClient(c KinesisConfig) (MessageService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.Region),
	)

	if err != nil {
		return nil, err
	}

	return &KinesisClient{
		service: kinesis.NewFromConfig(cfg),
		config:  c,
	}, nil
}

func (c *KinesisClient) SendMessage(message string) error {
	_, err := c.service.PutRecord(context.TODO(), &kinesis.PutRecordInput{
		StreamName:   &c.config.StreamName,
		PartitionKey: aws.String("partitionKey"),
		Data:         []byte(message),
	})

	return err
}

func (c *KinesisClient) ReadMessage() (*Message, error) {
	result, err := c.service.GetShardIterator(context.TODO(), &kinesis.GetShardIteratorInput{
		StreamName:        &c.config.StreamName,
		ShardId:           &c.config.ShardId,
		ShardIteratorType: types.ShardIteratorTypeTrimHorizon,
	})

	if err != nil {
		return nil, err
	}

	out, err := c.service.GetRecords(context.TODO(), &kinesis.GetRecordsInput{
		ShardIterator: result.ShardIterator,
		Limit:         aws.Int32(1),
	})

	if err != nil {
		return nil, err
	}

	if len(out.Records) == 0 {
		return nil, nil
	}

	return &Message{Body: string(out.Records[0].Data)}, nil
}
