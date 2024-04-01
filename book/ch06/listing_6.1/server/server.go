package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	shipping "listing_6.1"
)

type server struct {
	shipping.UnimplementedShippingServiceServer
}

func (s *server) Create(ctx context.Context, in *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	time.Sleep(2 * time.Second) // simulated delay
	return &shipping.CreateShippingResponse{ShippingId: 1243}, nil
}

func main() {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	go func() {
		res := randomFunc("a")
		log.Println(res)
		cancel()
	}()
	go func() {
		res := randomFunc("b")
		log.Println(res)
		cancel()
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	shipping.RegisterShippingServiceServer(grpcServer, &server{})
	grpcServer.Serve(listener)
}

func randomFunc(name string) string {
	min := 3
	max := 7
	sleepTime := rand.Intn(max-min+1) + min
	time.Sleep(time.Duration(sleepTime) * 1000000)
	return "hello from " + name
}
