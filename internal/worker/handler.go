package worker

import (
	"bytes"
	"encoding/json"
	"github.com/sony/gobreaker"
	"go-ambassador-request/internal/acceptor"
	"go-ambassador-request/internal/checker"
	"go-ambassador-request/pkg/config"
	"go-ambassador-request/pkg/notification"
	"io"
	"log"
	"net/http"
	"time"
)

type BackgroundWorker struct {
	repo *checker.DynamoDbRepository
	msg  *acceptor.SQSMessageSender
	cb   *gobreaker.CircuitBreaker
}

func NewBackgroundWorker(cfg *config.Config) *BackgroundWorker {
	repository, err := checker.NewDynamoDbRepository(cfg)
	if err != nil {
		log.Fatalf("Failed to create DynamoDB repo: %s", err)
	}
	message, err := acceptor.NewSQSMessageSender(cfg)
	if err != nil {
		log.Fatalf("Failed to create SQS msg: %s", err)
	}

	/*var st gobreaker.Settings
	st.Name = "AmbassadorRequest"
	st.Interval = time.Duration(30) * time.Second
	st.Timeout = time.Duration(60) * time.Second
	st.MaxRequests = 5
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 3 && failureRatio >= 0.6
	}*/

	return &BackgroundWorker{
		repo: repository,
		msg:  message,
		// cb:   gobreaker.NewCircuitBreaker(st),
	}
}

func (bw BackgroundWorker) RunWorker() {
	message := make(chan *notification.Message)
	errors := make(chan error)

	go func() {
		for {
			bw.msg.ReadMessage(message, errors)
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		select {
		case msg := <-message:
			go bw.processSQSMessage(msg)
		case err := <-errors:
			log.Println("Erro:", err)
		}
	}
}

func (bw BackgroundWorker) processSQSMessage(message *notification.Message) {
	log.Println("Received message:", message.ID)

	var request acceptor.Request
	errUnmarshal := json.Unmarshal([]byte(message.Body), &request.Data)
	if errUnmarshal != nil {
		log.Println("error converting message:", request.Data)
		return
	}

	err := bw.repo.PutItem(message.Body)
	if err != nil {
		log.Println("Error writing to DynamoDB:", err)
		return
	}

	req, _ := http.NewRequest(http.MethodPost, request.MetadataURL(), bytes.NewBuffer([]byte(message.Body)))

	client := &http.Client{}
	res, err := client.Do(req)

	body, readAllError := io.ReadAll(res.Body)
	if readAllError != nil {
		log.Println("Error converting response:", readAllError)
		return
	}

	response := Response{
		Body:       string(body),
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}

	request.AddMetadataAttempts(message.ApproximateReceiveCount)
	request.AddMetadataResponse(response)

	dataJson, _ := json.Marshal(request.Data)
	bw.repo.UpdateItem(string(dataJson))

	if err != nil || res.StatusCode != http.StatusOK {
		log.Println("Message not processed, a new attempt will be made")
	} else {
		bw.msg.DeleteMessage(message.ReceiptHandle)
		log.Println("Processed message")
	}
}
