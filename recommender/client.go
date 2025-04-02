package recommender

import (
	"context"
	"github.com/rasadov/EcommerceMicroservices/recommender/generated/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.RecommenderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewRecommenderServiceClient(conn)
	return &Client{conn, client}, nil
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) GetRecommendation(ctx context.Context, userId string) (*pb.RecommendationResponse, error) {
	return client.service.GetRecommendations(
		ctx,
		&pb.RecommendationRequest{UserId: userId},
	)
}
