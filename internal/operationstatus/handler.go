package operationstatus

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Handler struct {
	sender *DynamoDbRepository
}

func (h *Handler) OperationStatusChecker(context *gin.Context) {
	panic("Not implemented")
}

func NewHandler() *Handler {
	repository, err := NewDynamoDbRepository()
	if err != nil {
		log.Fatalf("Failed to create DynamoDB repository: %s", err)
	}
	return &Handler{sender: repository}
}
