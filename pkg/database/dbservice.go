package database

import (
	"errors"
)

type DbService interface {
	GetItem(id string) (any, error)
	PutItem(item interface{}) error
}

type DbType int

const (
	DynamoDB DbType = iota
	DocumentDB
)

func NewDatabaseService(databaseType DbType, config interface{}) (DbService, error) {
	switch databaseType {
	case DynamoDB:
		c, ok := config.(DynamoDBConfig)
		if !ok {
			return nil, errors.New("invalid config for DynamoDB")
		}
		return NewDynamoDBClient(c)
	case DocumentDB:
		c, ok := config.(DocumentDBConfig)
		if !ok {
			return nil, errors.New("invalid config for DocumentDB")
		}
		return NewDocumentDBClient(c)
	default:
		return nil, errors.New("unknown database type")
	}
}
