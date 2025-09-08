package graph

import (
	"context"
	"log"
	"time"
)

type accountResolver struct {
	server *Server
}

func (resolver *accountResolver) ID(ctx context.Context, obj *Account) (int, error) {
	return int(obj.ID), nil
}

func (resolver *accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	orderList, err := resolver.server.orderClient.GetOrdersForAccount(ctx, obj.ID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var orders []*Order
	for _, order := range orderList {
		var products []*OrderedProduct
		for _, orderedProduct := range order.Products {
			products = append(products, &OrderedProduct{
				ID:          orderedProduct.ID,
				Name:        orderedProduct.Name,
				Description: orderedProduct.Description,
				Price:       orderedProduct.Price,
				Quantity:    int(orderedProduct.Quantity),
			})
		}
		orders = append(orders, &Order{
			ID:         int(order.ID),
			CreatedAt:  order.CreatedAt,
			TotalPrice: order.TotalPrice,
			Products:   products,
		})
	}

	return orders, nil
}
