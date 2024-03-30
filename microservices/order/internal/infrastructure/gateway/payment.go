package gateway

import (
	"context"

	"github.com/peppydays/grpc-in-go-service/idl/contract/go/payment"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentGateway struct {
	client payment.PaymentClient
}

func NewPaymentGateway(paymentServiceUrl string) (*PaymentGateway, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := payment.NewPaymentClient(conn)
	return &PaymentGateway{client: client}, nil
}

func (p *PaymentGateway) Charge(order *domain.Order) error {
	_, err := p.client.Create(
		context.Background(),
		&payment.CreatePaymentRequest{
			UserId:     order.CustomerID,
			OrderId:    order.ID,
			TotalPrice: order.TotalPrice(),
		},
	)
	return err
}
