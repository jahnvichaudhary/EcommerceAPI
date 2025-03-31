package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rasadov/EcommerceMicroservices/account"
	"github.com/rasadov/EcommerceMicroservices/order"
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

func (resolver *mutationResolver) Register(ctx context.Context, in RegisterInput) (*AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	token, err := resolver.server.accountClient.Register(ctx, in.Name, in.Email, in.Password)
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

func (resolver *mutationResolver) Login(ctx context.Context, in LoginInput) (*AuthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	token, err := resolver.server.accountClient.Login(ctx, in.Email, in.Password)
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

func (resolver *mutationResolver) CreateProduct(ctx context.Context, in CreateProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	accountId, err := account.GetUserIdInt(ctx, true)
	if err != nil {
		return nil, err
	}

	postProduct, err := resolver.server.productClient.PostProduct(ctx, in.Name, in.Description, in.Price, int64(accountId))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		ID:          postProduct.ID,
		Name:        postProduct.Name,
		Description: postProduct.Description,
		Price:       postProduct.Price,
		AccountID:   accountId,
	}, nil
}

func (resolver *mutationResolver) UpdateProduct(ctx context.Context, in UpdateProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	accountId, err := account.GetUserIdInt(ctx, true)
	if err != nil {
		return nil, err
	}

	updatedProduct, err := resolver.server.productClient.UpdateProduct(ctx, in.ID, in.Name, in.Description, in.Price, int64(accountId))
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          updatedProduct.ID,
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		AccountID:   accountId,
	}, nil
}

func (resolver *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	accountId, err := account.GetUserIdInt(ctx, true)
	if err != nil {
		return false, err
	}

	err = resolver.server.productClient.DeleteProduct(ctx, id, int64(accountId))
	if err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func (resolver *mutationResolver) CreateOrder(ctx context.Context, in OrderInput) (*Order, error) {
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

	o, err := resolver.server.orderClient.PostOrder(ctx, accountId, products)
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
