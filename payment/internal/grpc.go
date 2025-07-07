package internal

import (
	"context"
	"github.com/dodopayments/dodopayments-go"
	order "github.com/rasadov/EcommerceAPI/order/client"
	"github.com/rasadov/EcommerceAPI/payment/proto/pb"
)

type grpcServer struct {
	pb.UnimplementedPaymentServiceServer
	service     Service
	orderClient *order.Client
}

func (s *grpcServer) Checkout(ctx context.Context, request *pb.CheckoutRequest) (*pb.RedirectResponse, error) {
	customer, err := s.service.FindOrCreateCustomer(ctx, request.UserId, request.Email, request.Name)
	if err != nil {
		return nil, err
	}
	currency := dodopayments.Currency(request.Currency)

	checkoutUrl, productId, err := s.service.GetCheckoutURL(ctx, request.Email, request.Name, request.RedirectURL, request.Price, currency)

	// We will use these transaction on webhooks
	err = s.service.RegisterTransaction(ctx, request.OrderId, request.UserId, request.Price, currency, customer.CustomerId, productId)

	if err != nil {
		return nil, err
	}

	return &pb.RedirectResponse{
		Url: checkoutUrl,
	}, nil
}

func (s *grpcServer) GetCustomerPortal(ctx context.Context, request *pb.CustomerPortalRequest) (*pb.RedirectResponse, error) {
	customer, err := s.service.FindOrCreateCustomer(ctx, request.UserId, *request.Email, *request.Name)

	if err != nil {
		return nil, err
	}

	link, err := s.service.GetCustomerPortal(ctx, customer)

	if err != nil {
		return nil, err
	}
	return &pb.RedirectResponse{
		Url: link,
	}, nil
}
