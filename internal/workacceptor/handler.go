package workacceptor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
)

type Request struct {
	Data map[string]interface{} `json:"data"`
}

type Handler struct {
	sender   *SQSMessageSender
	hostname string
}

func NewHandler() *Handler {
	sqsSender, err := NewSQSMessageSender()
	if err != nil {
		log.Fatalf("Failed to create SQS sender: %s", err)
	}

	return &Handler{
		sender: sqsSender,
		//hostname: os.Getenv("API_HOSTNAME"),
		hostname: "localhost",
	}
}

func (h *Handler) ProcessingWorkAcceptor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed for this route", http.StatusMethodNotAllowed)
		return
	}

	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestID, _ := uuid.NewV4()
	rqs := fmt.Sprintf("http://%s/api/v1/checker/%s", h.hostname, requestID)

	request.Data["metadata"] = map[string]string{
		"id": requestID.String(),
	}

	payload, err := json.Marshal(request.Data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.sender.SendMessage(string(payload))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"message":     "Request Accepted for Processing",
		"ProxyStatus": rqs,
	})
}
