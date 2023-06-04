package main

import (
	"github.com/gin-gonic/gin"
	"go-ambassador-request/internal/operationstatus"
	"go-ambassador-request/internal/workacceptor"
	"log"
)

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.POST("/acceptor", workacceptor.NewHandler().ProcessingWorkAcceptor)
	v1.GET("/checker", operationstatus.NewHandler().OperationStatusChecker)

	// go processSQSMessages()

	err := router.Run()
	if err != nil {
		log.Fatal("Error run server:", err)
	}
}
