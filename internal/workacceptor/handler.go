package workacceptor

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Data map[string]interface{} `json:"data"`
}

type Handler struct {
	sender *SQSMessageSender
}

func NewHandler() *Handler {
	sqsSender, err := NewSQSMessageSender()
	if err != nil {
		log.Fatalf("Failed to create SQS sender: %s", err)
	}
	return &Handler{sender: sqsSender}
}

func (h *Handler) ProcessingWorkAcceptor(c *gin.Context) {
	var request Request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if request.Data == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	requestID, _ := uuid.NewV4()
	rqs := fmt.Sprintf("http://%s/api/v1/OperationStatusChecker/%s", os.Getenv("API_HOSTNAME"), requestID)

	addMetadata(&request, requestID)

	payload, err := json.Marshal(request.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error processing request"})
		return
	}

	err = h.sender.SendMessage(string(payload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error processing request"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":     "Request Accepted for Processing",
		"ProxyStatus": rqs,
	})
}

func addMetadata(request *Request, requestID uuid.UUID) {
	request.Data["metadata"] = map[string]string{
		"RequestID": requestID.String(),
	}
}
