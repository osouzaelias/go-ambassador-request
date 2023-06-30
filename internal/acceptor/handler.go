package acceptor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	var headers = map[string][]string{}
	for name, values := range r.Header {
		if name != "AR-Destination-Url" && strings.HasPrefix(name, "AR-") {
			headers[name] = values
		}
	}

	if len(headers) > 0 {
		request.Data["_headers"] = headers
	}

	requestID, _ := uuid.NewV4()
	proxyStatus := fmt.Sprintf("http://%s/api/v1/checker/%s", h.hostname, requestID)

	request.Data["_id"] = requestID.String()
	request.Data["_url"] = url

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
		"proxyStatus": proxyStatus,
	})
}
