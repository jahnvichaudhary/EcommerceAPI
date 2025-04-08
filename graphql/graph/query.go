package graph

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/rasadov/EcommerceAPI/pkg/auth"
)

type queryResolver struct {
	server *Server
}

func (resolver *queryResolver) Accounts(
	ctx context.Context,
	pagination *PaginationInput,
	id *string,
) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		res, err := resolver.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{{
			ID:    strconv.Itoa(int(res.ID)),
			Name:  res.Name,
			Email: res.Email,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}
	accountList, err := resolver.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var accounts []*Account
	for _, account := range accountList {
		account := &Account{
			ID:    strconv.Itoa(int(account.ID)),
			Name:  account.Name,
			Email: account.Email,
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (resolver *queryResolver) Product(
	ctx context.Context,
	pagination *PaginationInput,
	query, id *string,
	viewedProductsIds []*string,
	byAccountId *bool,
) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get single
	if id != nil {
		res, err := resolver.server.productClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Product{{
			ID:          res.ID,
			Name:        res.Name,
			Description: res.Description,
			Price:       res.Price,
		}}, nil
	}
	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	// Get recommendations
	if viewedProductsIds != nil {
		productIds := make([]string, len(viewedProductsIds))
		for i, id := range viewedProductsIds {
			productIds[i] = *id
		}
		res, err := resolver.server.recommenderClient.GetRecommendationBasedOnViewed(ctx, productIds, skip, take)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		productList := res.GetRecommendedProducts()
		var products []*Product
		for _, product := range productList {
			products = append(products,
				&Product{
					ID:          product.Id,
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
				},
			)
		}
		return products, nil
	}

	if byAccountId != nil && *byAccountId {
		accountId := auth.GetUserId(ctx, true)
		if accountId == "" {
			return nil, errors.New("unauthorized")
		}
		skip = 0
		take = 100
		res, err := resolver.server.recommenderClient.GetRecommendationForUser(ctx, accountId, skip, take)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		productList := res.GetRecommendedProducts()
		var products []*Product
		for _, product := range productList {
			products = append(products,
				&Product{
					ID:          product.Id,
					Name:        product.Name,
					Description: product.Description,
					Price:       product.Price,
				},
			)
		}
		return products, nil
	}

	q := ""
	if query != nil {
		q = *query
	}
	productList, err := resolver.server.productClient.GetProducts(ctx, skip, take, nil, q)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*Product
	for _, product := range productList {
		products = append(products,
			&Product{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			},
		)
	}

	return products, nil
}

func (pagination PaginationInput) bounds() (uint64, uint64) {
	skipValue := uint64(0)
	takeValue := uint64(100)
	if pagination.Skip != 0 {
		skipValue = uint64(pagination.Skip)
	}
	if pagination.Take != 100 {
		takeValue = uint64(pagination.Take)
	}
	return skipValue, takeValue
}
