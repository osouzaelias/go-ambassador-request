package main

import (
	"go-ambassador-request/internal/acceptor"
	"go-ambassador-request/internal/checker"
	"go-ambassador-request/internal/worker"
	"log"
	"net/http"
)

func main() {
	bgw := worker.NewBackgroundWorker()
	go bgw.RunWorker()

	http.HandleFunc("/api/v1/acceptor", acceptor.NewHandler().ProcessingWorkAcceptor)
	http.HandleFunc("/api/v1/checker/", checker.NewHandler().OperationStatusChecker)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
