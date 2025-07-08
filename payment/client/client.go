package client

import (
	"github.com/rasadov/EcommerceAPI/payment/proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.PaymentServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	C := pb.NewPaymentServiceClient(conn)
	return &Client{conn, C}, nil
}

func (client *Client) Close() {
	err := client.conn.Close()
	if err != nil {
		log.Println(err)
	}
}
