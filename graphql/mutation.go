package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rasadov/EcommerceMicroservices/account"
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

func (r *mutationResolver) Register(ctx context.Context, in RegisterInput) (*AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	token, err := r.server.accountClient.Register(ctx, in.Name, in.Email, in.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ginContext, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, errors.New("could not retrieve gin context")
	}
	ginContext.SetCookie("token", token, 3600, "/", "localhost", false, true)

	return &AuthResponse{Token: token}, nil
}

func (r *mutationResolver) Login(ctx context.Context, in LoginInput) (*AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	token, err := r.server.accountClient.Login(ctx, in.Email, in.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ginContext, ok := ctx.Value("GinContextKey").(*gin.Context)
	if !ok {
		return nil, errors.New("could not retrieve gin context")
	}
	ginContext.SetCookie("token", token, 3600, "/", "localhost", false, true)

	return &AuthResponse{Token: token}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	postProduct, err := r.server.productClient.PostProduct(ctx, in.Name, in.Description, in.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	accountId, err := account.GetUserIdInt(ctx, true)
	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          postProduct.ID,
		Name:        postProduct.Name,
		Description: postProduct.Description,
		Price:       product.StringToFloat(postProduct.Price),
		AccountID:   accountId,
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

	accountId := account.GetUserId(ctx, true)
	if accountId == "" {
		return nil, errors.New("unauthorized")
	}

	o, err := r.server.orderClient.PostOrder(ctx, accountId, products)
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
