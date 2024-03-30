package main

import (
	"fmt"
	"log"
	"net"

	"github.com/peppydays/grpc-in-go-service/microservices/order/configs"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/application"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/infrastructure/gateway"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/infrastructure/repository"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/interface/rpc"
)

func main() {
	orderRepository, err := repository.NewGORMOrderRepository(configs.GetDataSourceURL())
	if err != nil {
		log.Fatalf("failed to connect to the database due to error %v", err)
	}

	paymentGateway, err := gateway.NewPaymentGateway(configs.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf("failed to initialise payment stub due to error %v", err)
	}

	orderService := application.NewOrderService(orderRepository, paymentGateway)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", configs.GetApplicationPort()))
	if err != nil {
		log.Fatalf("failed to bind port 8080 due to error %v", err)
	}

	rpc.RunAPI(listener, orderService)
}
