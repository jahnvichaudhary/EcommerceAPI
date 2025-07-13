package internal

import (
	"context"

	"github.com/dodopayments/dodopayments-go"
	order "github.com/rasadov/EcommerceAPI/order/client"
	"github.com/rasadov/EcommerceAPI/payment/proto/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type grpcServer struct {
	pb.UnimplementedPaymentServiceServer
	service     Service
	orderClient *order.Client
}

func (s *grpcServer) Checkout(ctx context.Context, request *pb.CheckoutRequest) (*wrapperspb.StringValue, error) {
	customer, err := s.service.FindOrCreateCustomer(ctx, request.UserId, request.Email, request.Name)
	if err != nil {
		return nil, err
	}
	currency := dodopayments.Currency(request.Currency)

	checkoutUrl, productId, err := s.service.GetCheckoutURL(ctx, request.Email, request.Name, request.RedirectURL, request.PriceCents, currency)
	if err != nil {
		return nil, err
	}

	// We will use these transaction on webhooks
	err = s.service.RegisterTransaction(ctx, request.OrderId, request.UserId, request.PriceCents, currency, customer.CustomerId, productId)
	if err != nil {
		return nil, err
	}

	return &wrapperspb.StringValue{
		Value: checkoutUrl,
	}, nil
}

func (s *grpcServer) CreateCustomerPortalSession(ctx context.Context, request *pb.CustomerPortalRequest) (*wrapperspb.StringValue, error) {
	customer, err := s.service.FindOrCreateCustomer(ctx, request.UserId, *request.Email, *request.Name)

	if err != nil {
		return nil, err
	}

	link, err := s.service.CreateCustomerPortalSession(ctx, customer)

	if err != nil {
		return nil, err
	}
	return &wrapperspb.StringValue{
		Value: link,
	}, nil
}
