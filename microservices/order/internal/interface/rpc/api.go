package api

import (
	"context"
	"log"
	"net"

	"github.com/peppydays/grpc-microservices-in-go/idl/contract/go/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/peppydays/grpc-microservices-in-go/microservices/order/internal/domain"
)

type API struct {
	service domain.OrderService
	order.UnimplementedOrderServer
}

func RunAPI(listener net.Listener, service domain.OrderService) {
	api := API{service: service}
	grpcServer := grpc.NewServer()
	order.RegisterOrderServer(grpcServer, api)
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server")
	}
}

func (api API) Create(context context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem

	for _, orderItem := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	newOrder := domain.NewOrder(request.UserId, orderItems)

	result, err := api.service.PlaceOrder(newOrder)
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{OrderId: result.ID}, nil
}
