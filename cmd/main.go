package main

import (
	"go-ambassador-request/internal/operationstatus"
	"go-ambassador-request/internal/workacceptor"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/acceptor", workacceptor.NewHandler().ProcessingWorkAcceptor)
	http.HandleFunc("/api/v1/checker/", operationstatus.NewHandler().OperationStatusChecker)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
