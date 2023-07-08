package checker

import (
	"go-ambassador-request/pkg/config"
	"go-ambassador-request/pkg/database"
)

type DynamoDbRepository struct {
	client database.DbService
}

func NewDynamoDbRepository(cfg *config.Config) (*DynamoDbRepository, error) {
	dbConfig := database.DynamoDBConfig{
		Region: cfg.Region,
		Table:  cfg.Table,
	}

	client, err := database.NewDatabaseService(database.DynamoDB, dbConfig)
	if err != nil {
		return nil, err
	}

	return &DynamoDbRepository{client: client}, nil
}

func (s *DynamoDbRepository) GetItem(id string) (interface{}, error) {
	return s.client.GetItem(id)
}

func (s *DynamoDbRepository) PutItem(msg string) error {
	return s.client.PutItem(msg)
}

func (s *DynamoDbRepository) UpdateItem(item string) error {
	return s.client.UpdateItem(item)
}
