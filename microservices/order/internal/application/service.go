package application

import (
	"strings"

	"github.com/peppydays/grpc-in-go-service/microservices/order/internal/domain"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderService struct {
	repository     domain.OrderRepository
	paymentGateway domain.PaymentGateway
}

func NewOrderService(repository domain.OrderRepository, paymentGateway domain.PaymentGateway) *OrderService {
	return &OrderService{
		repository:     repository,
		paymentGateway: paymentGateway,
	}
}

func (s OrderService) PlaceOrder(order domain.Order) (domain.Order, error) {
	err := s.repository.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	err = s.paymentGateway.Charge(&order)
	if err != nil {
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")

		badReq := &errdetails.BadRequest{}
		st := status.Convert(err)
		var allErrors []string
		for _, detail := range st.Details() {
			switch t := detail.(type) {
			case *errdetails.BadRequest:
				for _, violation := range t.GetFieldViolations() {
					allErrors = append(allErrors, violation.Description)
				}
			}
		}
		fieldErr := &errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: strings.Join(allErrors, "\n"),
		}
		badReq.FieldViolations = append(badReq.FieldViolations, fieldErr)

		statusWithDetails, _ := orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDetails.Err()
	}

	return order, nil
}
