package internal

import (
	"context"
	"fmt"
	"github.com/dodopayments/dodopayments-go"
	order "github.com/rasadov/EcommerceAPI/order/client"
	"github.com/rasadov/EcommerceAPI/payment/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type grpcServer struct {
	pb.UnimplementedPaymentServiceServer
	service     Service
	orderClient *order.Client
}

func ListenGRPC(service Service, orderURL string, port int) error {
	orderClient, err := order.NewClient(orderURL)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		orderClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterPaymentServiceServer(serv, &grpcServer{
		pb.UnimplementedPaymentServiceServer{},
		service,
		orderClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
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
