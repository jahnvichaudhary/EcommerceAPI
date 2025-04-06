package server

import (
	"context"
	"fmt"
	"github.com/rasadov/EcommerceAPI/account/internal/user"
	pb2 "github.com/rasadov/EcommerceAPI/account/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

type grpcServer struct {
	pb2.UnimplementedAccountServiceServer
	service user.Service
}

func ListenGRPC(service user.Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()

	pb2.RegisterAccountServiceServer(serv, &grpcServer{
		UnimplementedAccountServiceServer: pb2.UnimplementedAccountServiceServer{},
		service:                           service})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (server *grpcServer) Register(ctx context.Context, request *pb2.RegisterRequest) (*pb2.AuthResponse, error) {
	token, err := server.service.Register(ctx, request.Name, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &pb2.AuthResponse{
		Token: token,
	}, nil
}

func (server *grpcServer) Login(ctx context.Context, request *pb2.LoginRequest) (*pb2.AuthResponse, error) {
	token, err := server.service.Login(ctx, request.Email, request.Password)
	if err != nil {
		return nil, err
	}
	return &pb2.AuthResponse{
		Token: token,
	}, nil
}

func (server *grpcServer) GetAccount(ctx context.Context, r *pb2.GetAccountRequest) (*pb2.AccountResponse, error) {
	a, err := server.service.GetAccount(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb2.AccountResponse{Account: &pb2.Account{
		Id:   strconv.Itoa(int(a.ID)),
		Name: a.Name,
	}}, nil
}

func (server *grpcServer) GetAccounts(ctx context.Context, r *pb2.GetAccountsRequest) (*pb2.GetAccountsResponse, error) {
	res, err := server.service.GetAccounts(ctx, r.Skip, r.Take)
	if err != nil {
		return nil, err
	}
	var accounts []*pb2.Account
	for _, p := range res {
		accounts = append(accounts, &pb2.Account{
			Id:   strconv.Itoa(int(p.ID)),
			Name: p.Name,
		},
		)
	}
	return &pb2.GetAccountsResponse{Accounts: accounts}, nil
}
