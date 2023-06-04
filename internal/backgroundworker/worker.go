package backgroundworker

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go-ambassador-request/internal/workacceptor"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func processSQSMessages() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg)
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	queueURL := os.Getenv("SQS_QUEUE_URL")

	for {
		// Continuously poll the SQS queue
		output, err := sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: 10,
		})
		if err != nil {
			log.Println("Error receiving SQS message: ", err)
			time.Sleep(1 * time.Minute) // Sleep for a minute before trying again
			continue
		}

		for _, message := range output.Messages {
			go processMessage(message, sqsClient, dynamodbClient, queueURL) // Process each message in a separate goroutine
		}

		time.Sleep(1 * time.Minute) // Sleep for a minute before checking the queue again
	}
}

func processMessage(message types.Message, sqsClient *sqs.Client, dynamodbClient *dynamodb.Client, queueURL string) {
	// Process the message
	log.Println("Received message: ", *message.Body)

	var request workacceptor.Request
	errUnmarshal := json.Unmarshal([]byte(*message.Body), &request)
	if errUnmarshal != nil {
		log.Fatal("error convert message")
	}

	item, _ := attributevalue.MarshalMap(request)

	// Write the message to DynamoDB
	_, err := dynamodbClient.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String("your-dynamodb-table-name"),
		Item:      item,
	})
	if err != nil {
		log.Println("Error writing to DynamoDB: ", err)
		return
	}

	// Delete the message from the SQS queue
	_, err = sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      &queueURL,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		log.Println("Error deleting SQS message: ", err)
	}
}
