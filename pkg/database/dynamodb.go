package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strings"
)

type DynamoDBConfig struct {
	Region string
	Table  string
}

type DynamoDBClient struct {
	service *dynamodb.Client
	config  DynamoDBConfig
}

func NewDynamoDBClient(c DynamoDBConfig) (DbService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(c.Region))
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBClient{service: client, config: c}, nil
}

func (c *DynamoDBClient) GetItem(id string) (interface{}, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(c.config.Table),
		Key: map[string]types.AttributeValue{
			"_id": &types.AttributeValueMemberS{Value: id},
		},
	}

	out, err := c.service.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if out.Item == nil {
		return nil, nil
	}

	item := make(map[string]interface{})
	err = attributevalue.UnmarshalMap(out.Item, &item)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (c *DynamoDBClient) PutItem(jsonStr string) error {
	var item map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &item)
	if err != nil {
		fmt.Println("error unmarshalling JSON:", err)
		return err
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling map:", err)
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(c.config.Table),
		Item:      av,
	}

	_, err = c.service.PutItem(context.TODO(), input)
	return err
}

func (c *DynamoDBClient) UpdateItem(jsonStr string) error {
	var item map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &item)
	if err != nil {
		fmt.Println("error unmarshalling JSON:", err)
		return err
	}

	primaryKey := item["_id"].(string)
	delete(item, "_id")

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		fmt.Println("Got error marshalling map:", err)
		return err
	}

	updateExpression := "SET"
	expressionAttributeValues := map[string]types.AttributeValue{}
	expressionAttributeNames := map[string]string{}

	for key, value := range av {
		expressionAttributeValues[":"+key] = value
		expressionAttributeNames["#"+key] = key
		updateExpression += " #" + key + " = :" + key + ","
	}
	updateExpression = strings.TrimRight(updateExpression, ",")

	input := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(c.config.Table),
		Key:                       map[string]types.AttributeValue{"_id": &types.AttributeValueMemberS{Value: primaryKey}},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		ExpressionAttributeNames:  expressionAttributeNames,
		ReturnValues:              types.ReturnValueAllNew,
	}

	_, err = c.service.UpdateItem(context.TODO(), input)
	return err
}
