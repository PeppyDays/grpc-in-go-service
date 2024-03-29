package main

import (
	"fmt"
	"log"
	"net"

	"github.com/peppydays/grpc-in-go-service/microservices/order/configs"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/application"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/infrastructure/repository"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/interface/rpc"
)

func main() {
	repository, err := repository.NewGORMOrderRepository(configs.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to the database due to error %v", err)
	}
	service := application.NewOrderService(repository)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.GetApplicationPort()))
	if err != nil {
		log.Fatalf("failed to bind port 8080 due to error %v", err)
	}

	rpc.RunAPI(listener, service)
}
