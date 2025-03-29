package account

import (
	"context"
	"fmt"
	"github.com/rasadov/EcommerceMicroservices/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()

	pb.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb.UnimplementedAccountServiceServer{},
		service:                           s})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.AuthResponse, error) {
	token, err := s.service.Register(ctx, request.Name, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{
		Token: token,
	}, nil
}

func (s *grpcServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.AuthResponse, error) {
	token, err := s.service.Login(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &pb.AuthResponse{
		Token: token,
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, r *pb.GetAccountRequest) (*pb.AccountResponse, error) {
	a, err := s.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.AccountResponse{Account: &pb.Account{
		Id:   strconv.Itoa(int(a.ID)),
		Name: a.Name,
	}}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, r *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	res, err := s.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	var accounts []*pb.Account
	for _, p := range res {
		accounts = append(accounts, &pb.Account{
			Id:   strconv.Itoa(int(p.ID)),
			Name: p.Name,
		},
		)
	}
	return &pb.GetAccountsResponse{Accounts: accounts}, nil
}
