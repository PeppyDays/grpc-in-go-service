package repository

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/peppydays/grpc-in-go-service/microservices/payment/internal/domain"
)

type Payment struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderID    int64
	TotalPrice float32
}

type GORMPaymentRepository struct {
	db *gorm.DB
}

func (a GORMPaymentRepository) Get(ctx context.Context, id string) (domain.Payment, error) {
	var paymentDataModel Payment
	res := a.db.WithContext(ctx).First(&paymentDataModel, id)
	payment := domain.Payment{
		ID:         int64(paymentDataModel.ID),
		CustomerID: paymentDataModel.CustomerID,
		Status:     paymentDataModel.Status,
		OrderID:    paymentDataModel.OrderID,
		TotalPrice: paymentDataModel.TotalPrice,
		CreatedAt:  paymentDataModel.CreatedAt.UnixNano(),
	}
	return payment, res.Error
}

func (a GORMPaymentRepository) Save(ctx context.Context, payment *domain.Payment) error {
	paymentDataModel := Payment{
		CustomerID: payment.CustomerID,
		Status:     payment.Status,
		OrderID:    payment.OrderID,
		TotalPrice: payment.TotalPrice,
	}
	res := a.db.WithContext(ctx).Create(&paymentDataModel)
	if res.Error == nil {
		payment.ID = int64(paymentDataModel.ID)
	}
	return res.Error
}

func NewGORMPaymentRepository(dataSourceUrl string) (*GORMPaymentRepository, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceUrl), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	err := db.AutoMigrate(&Payment{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}
	return &GORMPaymentRepository{db: db}, nil
}
