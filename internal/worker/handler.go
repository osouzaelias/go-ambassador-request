package worker

import (
	"encoding/json"
	"fmt"
	"go-ambassador-request/internal/acceptor"
	"go-ambassador-request/internal/checker"
	"go-ambassador-request/pkg/notification"
	"log"
	"time"
)

type BackgroundWorker struct {
	repo *checker.DynamoDbRepository
	msg  *acceptor.SQSMessageSender
}

func NewBackgroundWorker() *BackgroundWorker {
	repository, err := checker.NewDynamoDbRepository()
	if err != nil {
		log.Fatalf("Failed to create DynamoDB repo: %s", err)
	}
	message, err := acceptor.NewSQSMessageSender()
	if err != nil {
		log.Fatalf("Failed to create SQS msg: %s", err)
	}
	return &BackgroundWorker{repo: repository, msg: message}
}

func (bw BackgroundWorker) RunWorker() {
	messages := make(chan *notification.Message)
	errors := make(chan error)

	go func() {
		for {
			bw.msg.ReadMessage(messages, errors)
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		select {
		case msg := <-messages:
			bw.processSQSMessage(msg)
		case err := <-errors:
			fmt.Println("Erro:", err)
		}
	}
}

func (bw BackgroundWorker) processSQSMessage(message *notification.Message) {
	// Process the message
	log.Println("Received message: ", message.ID)

	var request acceptor.Request
	errUnmarshal := json.Unmarshal([]byte(message.Body), &request)
	if errUnmarshal != nil {
		log.Fatal("error convert message")
	}

	err := bw.repo.PutItem(message.Body)

	if err != nil {
		log.Println("Error writing to DynamoDB: ", err)
		return
	}

	err = bw.msg.DeleteMessage(message.ReceiptHandle)

	if err != nil {
		log.Println("Error deleting SQS message: ", err)
	}
}
