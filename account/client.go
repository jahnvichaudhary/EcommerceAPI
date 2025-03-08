package account

import (
	"context"
	"github.com/rasadov/EcommerceMicroservices/account/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	C := pb.NewAccountServiceClient(conn)
	return &Client{conn, C}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.service.PostAccount(
		ctx,
		&pb.PostAccountRequest{Name: name},
	)

	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.GetId(),
		Name: r.Account.GetName(),
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, Id string) (*Account, error) {
	r, err := c.service.GetAccount(
		ctx,
		&pb.GetAccountRequest{Id: Id},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.GetId(),
		Name: r.Account.GetName(),
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	r, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{Take: take, Skip: skip},
	)
	if err != nil {
		return nil, err
	}
	var accounts []Account
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   a.GetId(),
			Name: a.GetName(),
		})
	}
	return accounts, nil
}
