package client

import (
	"context"
	"github.com/rasadov/EcommerceAPI/account/internal/user"
	pb2 "github.com/rasadov/EcommerceAPI/account/proto/pb"
	"github.com/rasadov/EcommerceAPI/pkg/auth"

	"google.golang.org/grpc"
	"strconv"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb2.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	C := pb2.NewAccountServiceClient(conn)
	return &Client{conn, C}, nil
}

func NewJwtService(secretKey, issuer string) auth.AuthService {
	return &auth.JwtService{
		SecretKey: secretKey,
		Issuer:    issuer,
	}
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) Register(ctx context.Context, name, email, password string) (string, error) {
	response, err := client.service.RegisterAccount(ctx, &pb2.RegisterRequest{
		Name:     name,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return response.Token, nil
}

func (client *Client) Login(ctx context.Context, email, password string) (string, error) {
	response, err := client.service.LoginAccount(ctx, &pb2.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	return response.Token, nil
}

func (client *Client) GetAccount(ctx context.Context, Id string) (*user.Account, error) {
	r, err := client.service.GetAccount(
		ctx,
		&pb2.GetAccountRequest{Id: Id},
	)
	if err != nil {
		return nil, err
	}
	accountId, _ := strconv.ParseInt(r.Account.GetId(), 10, 64)
	return &user.Account{
		ID:    uint(accountId),
		Name:  r.Account.GetName(),
		Email: r.Account.GetEmail(),
	}, nil
}

func (client *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]user.Account, error) {
	r, err := client.service.GetAccounts(
		ctx,
		&pb2.GetAccountsRequest{Take: take, Skip: skip},
	)
	if err != nil {
		return nil, err
	}
	var accounts []user.Account
	for _, a := range r.Accounts {
		accountId, _ := strconv.ParseInt(a.GetId(), 10, 64)
		accounts = append(accounts, user.Account{
			ID:    uint(accountId),
			Name:  a.GetName(),
			Email: a.GetEmail(),
		})
	}
	return accounts, nil
}
