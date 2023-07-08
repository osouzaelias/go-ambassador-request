package database

import (
	"errors"
)

type DbService interface {
	GetItem(id string) (any, error)
	PutItem(item string) error
	UpdateItem(item string) error
}

type DbType int

const (
	DynamoDB DbType = iota
)

func NewDatabaseService(databaseType DbType, config interface{}) (DbService, error) {
	switch databaseType {
	case DynamoDB:
		c, ok := config.(DynamoDBConfig)
		if !ok {
			return nil, errors.New("invalid config for DynamoDB")
		}
		return NewDynamoDBClient(c)
	default:
		return nil, errors.New("unknown database type")
	}
}
