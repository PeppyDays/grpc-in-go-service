package rpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/peppydays/grpc-in-go-service/idl/contract/go/payment"
	"github.com/peppydays/grpc-in-go-service/microservices/payment/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type API struct {
	service domain.PaymentService
	payment.UnimplementedPaymentServer
}

func RunAPI(listener net.Listener, service domain.PaymentService) {
	api := API{service: service}
	grpcServer := grpc.NewServer()
	payment.RegisterPaymentServer(grpcServer, api)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server")
	}
}

func (a API) Create(ctx context.Context, request *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	newPayment := domain.NewPayment(request.UserId, request.OrderId, request.TotalPrice)
	result, err := a.service.Charge(ctx, newPayment)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to charge due to %v", err)).Err()
	}

	return &payment.CreatePaymentResponse{PaymentId: result.ID}, nil
}
