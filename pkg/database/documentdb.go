package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocumentDBClient struct {
	collection *mongo.Collection
}

type DocumentDBConfig struct {
	Uri        string
	Database   string
	Collection string
}

func NewDocumentDBClient(c DocumentDBConfig) (DbService, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c.Uri))
	if err != nil {
		return nil, err
	}
	collection := client.Database(c.Database).Collection(c.Collection)
	return &DocumentDBClient{collection: collection}, nil
}

func (c *DocumentDBClient) GetItem(id string) (interface{}, error) {
	// Note: Simplified for brevity, actual implementation will depend on your data model
	// You would typically use a filter to find the document by its id
	var result interface{}
	err := c.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)
	return result, err
}

func (c *DocumentDBClient) PutItem(item interface{}) error {
	// Note: Simplified for brevity, actual implementation will depend on your data model
	// You would typically marshal the item into a bson.M
	_, err := c.collection.InsertOne(context.TODO(), item)
	return err
}
