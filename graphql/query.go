package main

import (
	"context"
	"github.com/rasadov/EcommerceMicroservices/product"
	"log"
	"strconv"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		res, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{{
			ID:   strconv.Itoa(int(res.ID)),
			Name: res.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accounts []*Account
	for _, a := range accountList {
		account := &Account{
			ID:   strconv.Itoa(int(a.ID)),
			Name: a.Name,
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *queryResolver) Product(ctx context.Context, pagination *PaginationInput, query, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		r, err := r.server.productClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Price:       product.StringToFloat(r.Price),
		}}, nil
	}
	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}
	productList, err := r.server.productClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product
	for _, a := range productList {
		products = append(products,
			&Product{
				ID:          a.ID,
				Name:        a.Name,
				Description: a.Description,
				Price:       product.StringToFloat(a.Price),
			},
		)
	}

	return products, nil
}

func (p PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if p.Skip != 0 {
		skipValue = uint64(p.Skip)
	}
	if p.Take != 100 {
		takeValue = uint64(p.Take)
	}
	return skipValue, takeValue
}
