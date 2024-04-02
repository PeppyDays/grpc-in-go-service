package e2e

import (
	"context"
	"log"
	"testing"

	"github.com/peppydays/grpc-in-go-service/idl/contract/go/order"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CreateOrderTestSuite struct {
	suite.Suite
	compose tc.ComposeStack
}

func (c *CreateOrderTestSuite) SetupSuite() {
	compose, err := tc.NewDockerCompose("docker-compose.yaml")
	if err != nil {
		log.Fatalf("Could not run compose stack: %v", err)
	}

	c.compose = compose
}

func (c *CreateOrderTestSuite) TearDownSuite() {
	c.compose.Down(context.Background(), tc.RemoveOrphans(true))
}

func (c *CreateOrderTestSuite) Test_Should_Create_Order() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("localhost:8080", opts...)
	if err != nil {
		log.Fatalf("Failed to connect order service. Err: %v", err)
	}

	defer conn.Close()

	orderClient := order.NewOrderClient(conn)
	createOrderResponse, err := orderClient.Create(context.Background(), &order.CreateOrderRequest{
		UserId: 23,
		OrderItems: []*order.OrderItem{
			{
				ProductCode: "CAM123",
				Quantity:    3,
				UnitPrice:   1.23,
			},
		},
	})
	c.Nil(err)

	getOrderResponse, err := orderClient.Get(context.Background(), &order.GetOrderRequest{OrderId: createOrderResponse.OrderId})
	c.Nil(err)
	c.Equal(int64(23), getOrderResponse.UserId)
	orderItem := getOrderResponse.OrderItems[0]
	c.Equal(float32(1.23), orderItem.UnitPrice)
	c.Equal(int32(3), orderItem.Quantity)
	c.Equal("CAM123", orderItem.ProductCode)
}

func TestCreateOrderTestSuite(t *testing.T) {
	suite.Run(t, new(CreateOrderTestSuite))
}
