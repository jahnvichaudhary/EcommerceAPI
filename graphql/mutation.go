package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/rasadov/EcommerceMicroservices/order"
	"github.com/rasadov/EcommerceMicroservices/product"
	"log"
	"strconv"
	"time"
)

var (
	ErrInvalidParameter = errors.New("invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	a, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Account{
		ID:   strconv.Itoa(int(a.ID)),
		Name: a.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	p, err := r.server.productClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       product.StringToFloat(p.Price),
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var products []order.OrderedProduct
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		u, err := strconv.ParseUint(p.ID, 10, 32)
		if err != nil {
			fmt.Println("Error:", err)
			return nil, err
		}
		products = append(products, order.OrderedProduct{
			ID:       uint(u),
			Quantity: uint32(p.Quantity),
		})
	}
	o, err := r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Order{
		ID:         strconv.Itoa(int(o.ID)),
		CreatedAt:  o.CreatedAt,
		TotalPrice: o.TotalPrice,
	}, nil
}
