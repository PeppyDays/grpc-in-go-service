package repository

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type GORMOrderRepositoryTestSuite struct {
	suite.Suite
	DataSourceUrl string
}

func (o *GORMOrderRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()
	port := "3306/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("root:welcome@tcp(%s:%s)/order?parseTime=true", host, port.Port())
	}
	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "welcome",
			"MYSQL_DATABASE":      "order",
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", dbURL).WithStartupTimeout(time.Second * 10).WithQuery("SELECT 1"),
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal("Failed to start MySQL database", err)
	}
	endpoint, _ := mysqlContainer.Endpoint(ctx, "")
	o.DataSourceUrl = fmt.Sprintf("root:welcome@tcp(%s)/order?parseTime=true", endpoint)
}

func (o *GORMOrderRepositoryTestSuite) Test_Should_Save_Order() {
	repository, err := NewGORMOrderRepository(o.DataSourceUrl)
	o.Nil(err)

	err = repository.Save(&domain.Order{})
	o.Nil(err)
}

func (o *GORMOrderRepositoryTestSuite) Test_Should_Get_Order() {
	repository, _ := NewGORMOrderRepository(o.DataSourceUrl)
	order := domain.NewOrder(2, []domain.OrderItem{
		{
			ProductCode: "CAM",
			Quantity:    5,
			UnitPrice:   1.32,
		},
	})

	_ = repository.Save(&order)
	loadedOrder, _ := repository.Get(order.ID)

	o.Equal(order.ID, loadedOrder.ID)
	o.Equal(order.CustomerID, loadedOrder.CustomerID)
	o.Equal(order.OrderItems, loadedOrder.OrderItems)
	o.Equal(order.Status, loadedOrder.Status)
}

func TestOrderDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(GORMOrderRepositoryTestSuite))
}
