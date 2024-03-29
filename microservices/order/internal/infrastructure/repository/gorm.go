package repository

import (
	"fmt"

	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type GORMOrderRepository struct {
	db *gorm.DB
}

func NewGORMOrderRepository(dataSourceURL string) (*GORMOrderRepository, error) {
	db, err := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection error: %v", err)
	}

	err = db.AutoMigrate(&Order{}, &OrderItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &GORMOrderRepository{db: db}, nil
}

func (r GORMOrderRepository) Get(id string) (domain.Order, error) {
	var orderDataModel Order
	res := r.db.First(&orderDataModel, id)

	if res.Error != nil {
		return domain.Order{}, res.Error
	}

	var orderItems []domain.OrderItem
	for _, orderItemDataModel := range orderDataModel.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItemDataModel.ProductCode,
			UnitPrice:   orderItemDataModel.UnitPrice,
			Quantity:    orderItemDataModel.Quantity,
		})
	}

	order := domain.Order{
		ID:         int64(orderDataModel.ID),
		CustomerID: orderDataModel.CustomerID,
		Status:     orderDataModel.Status,
		OrderItems: orderItems,
		CreatedAt:  orderDataModel.CreatedAt.UnixNano(),
	}

	return order, nil
}

func (r GORMOrderRepository) Save(order *domain.Order) error {
	var orderItemsDataModel []OrderItem
	for _, orderItem := range order.OrderItems {
		orderItemsDataModel = append(orderItemsDataModel, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	orderDataModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItemsDataModel,
	}

	res := r.db.Create(&orderDataModel)
	if res.Error == nil {
		order.ID = int64(orderDataModel.ID)
	}

	return res.Error
}
