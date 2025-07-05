package internal

import (
	"context"
	"fmt"
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
	return nil, nil
}

func (s *grpcServer) GetCustomerPortal(context.Context, *pb.CustomerPortalRequest) (*pb.RedirectResponse, error) {
	return nil, nil
}

func (s *grpcServer) FindOrCreateCustomer(context.Context, *pb.CustomerInput) (*pb.Customer, error) {
	return nil, nil
}
