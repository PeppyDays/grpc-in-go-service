package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"
	order "listing_6.5"
)

type server struct {
	order.UnimplementedOrderServiceServer
}

func (s *server) Create(ctx context.Context, in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var err error
	if rand.Intn(2) == 1 {
		err = errors.New("create order error")
	}
	return &order.CreateOrderResponse{OrderId: 1243}, err
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	order.RegisterOrderServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}
