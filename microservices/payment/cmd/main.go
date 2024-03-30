package main

import (
	"fmt"
	"log"
	"net"

	"github.com/peppydays/grpc-in-go-service/microservices/payment/configs"
	"github.com/peppydays/grpc-in-go-service/microservices/payment/internal/application"
	"github.com/peppydays/grpc-in-go-service/microservices/payment/internal/infrastructure/repository"
	"github.com/peppydays/grpc-in-go-service/microservices/payment/internal/interface/rpc"
)

func main() {
	paymentRepository, err := repository.NewGORMPaymentRepository(configs.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to the database due to error %v", err)
	}

	paymentService := application.NewPaymentService(paymentRepository)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.GetApplicationPort()))
	if err != nil {
		log.Fatalf("failed to bind port %d due to error %v", configs.GetApplicationPort(), err)
	}

	rpc.RunAPI(listener, paymentService)
}
