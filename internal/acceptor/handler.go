package acceptor

import (
	"encoding/json"
	"fmt"
	"go-ambassador-request/pkg/config"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofrs/uuid"
)

type Handler struct {
	sender   *SQSMessageSender
	hostname string
}

func NewHandler(cfg *config.Config) *Handler {
	sqsSender, err := NewSQSMessageSender(cfg)
	if err != nil {
		log.Fatalf("Failed to create SQS sender: %s", err)
	}

	return &Handler{
		sender:   sqsSender,
		hostname: os.Getenv("API_HOSTNAME"),
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

	url := r.Header.Get("AR-Destination-Url")
	if len(strings.TrimSpace(url)) == 0 {
		http.Error(w, "AR-Destination-Url header is missing", http.StatusBadRequest)
		return
	}

	requestID, _ := uuid.NewV4()
	proxyStatus := fmt.Sprintf("http://%s/api/v1/checker/%s", h.hostname, requestID)

	request.Data[metadataID] = requestID.String()
	request.Data[metadataURL] = url

	if err := h.sender.SendMessage(request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"message":     "Request Accepted for Processing",
		"proxyStatus": proxyStatus,
	})
}
