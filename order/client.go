package order

import (
	"context"
	"github.com/rasadov/EcommerceMicroservices/order/pb"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (client *Client) Close() {
	client.conn.Close()
}

func (client *Client) PostOrder(
	ctx context.Context,
	accountID string,
	products []OrderedProduct,
) (*Order, error) {
	var protoProducts []*pb.OrderProduct
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.OrderProduct{
			Id:       p.ID,
			Quantity: p.Quantity,
		})
	}
	r, err := client.service.PostOrder(
		ctx,
		&pb.PostOrderRequest{
			AccountId: accountID,
			Products:  protoProducts,
		},
	)
	if err != nil {
		return nil, err
	}

	// Create response order
	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)

	return &Order{
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountID:  newOrder.AccountId,
		Products:   products,
	}, nil
}

func (client *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	r, err := client.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Create response orders
	var orders []Order
	for _, orderProto := range r.Orders {
		orderId, _ := strconv.ParseInt(orderProto.Id, 10, 64)
		newOrder := Order{
			ID:         uint(orderId),
			TotalPrice: orderProto.TotalPrice,
			AccountID:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)

		var products []OrderedProduct
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}
	return orders, nil
}
