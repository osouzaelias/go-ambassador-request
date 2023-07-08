package main

import (
	"go-ambassador-request/internal/acceptor"
	"go-ambassador-request/internal/checker"
	"go-ambassador-request/internal/worker"
	"go-ambassador-request/pkg/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfig()

	bgw := worker.NewBackgroundWorker(cfg)
	go bgw.RunWorker()

	http.HandleFunc("/api/v1/acceptor", acceptor.NewHandler(cfg).ProcessingWorkAcceptor)
	http.HandleFunc("/api/v1/checker/", checker.NewHandler(cfg).OperationStatusChecker)
	log.Fatal(http.ListenAndServe(":8888", nil))
}
