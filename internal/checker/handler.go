package checker

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	repo *DynamoDbRepository
}

func NewHandler() *Handler {
	repository, err := NewDynamoDbRepository()
	if err != nil {
		log.Fatalf("Failed to create DynamoDB repo: %s", err)
	}
	return &Handler{repo: repository}
}

func (h *Handler) OperationStatusChecker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed for this route", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	idStr := parts[len(parts)-1]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v is an invalid id", id), http.StatusBadRequest)
		return
	}

	item, errItem := h.repo.GetItem(id.String())
	if errItem != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if item == nil {
		item = []byte{}
		w.WriteHeader(http.StatusNotFound)
	}

	jsonStr, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonStr)
}
