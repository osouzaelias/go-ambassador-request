package operationstatus

import (
	"go-ambassador-request/pkg/database"
	"os"
)

type DynamoDbRepository struct {
	client database.DbService
}

func NewDynamoDbRepository() (*DynamoDbRepository, error) {
	config := database.DynamoDBConfig{
		Region: os.Getenv("REGION"),
		Table:  os.Getenv("TABLE"),
	}

	client, err := database.NewDatabaseService(database.DynamoDB, config)
	if err != nil {
		return nil, err
	}

	return &DynamoDbRepository{client: client}, nil
}

func (s *DynamoDbRepository) GetItem(id string) (interface{}, error) {
	return s.client.GetItem(id)
}
