package product

import (
	"context"
	"errors"
)

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64, accountId int) (*Product, error)
	GetProduct(ctx context.Context, id string) (*Product, error)
	GetProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	GetProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
	UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId int) (*Product, error)
	DeleteProduct(ctx context.Context, productId string, accountId int) error
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	AccountID   int    `json:"accountID"`
}

type productService struct {
	repo Repository
}

func NewProductService(repository Repository) Service {
	return &productService{repository}
}

func (p productService) PostProduct(ctx context.Context, name, description string, price float64, accountId int) (*Product, error) {
	product := Product{
		Name:        name,
		Description: description,
		Price:       FloatToString(price),
		AccountID:   accountId,
	}

	err := p.repo.PutProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p productService) GetProduct(ctx context.Context, id string) (*Product, error) {
	return p.repo.GetProductById(ctx, id)
}

func (p productService) GetProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	return p.repo.ListProducts(ctx, skip, take)
}

func (p productService) GetProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	return p.repo.ListProductsWithIDs(ctx, ids)
}

func (p productService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	return p.repo.SearchProducts(ctx, query, skip, take)
}

func (p productService) UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId int) (*Product, error) {
	product, err := p.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}
	if product.AccountID != accountId {
		return nil, errors.New("unauthorized")
	}

	updatedProduct := Product{
		id,
		name,
		description,
		FloatToString(price),
		accountId,
	}
	err = p.repo.UpdateProduct(ctx, updatedProduct)
	if err != nil {
		return nil, err
	}
	return &updatedProduct, nil
}
func (p productService) DeleteProduct(ctx context.Context, productId string, accountId int) error {
	product, err := p.repo.GetProductById(ctx, productId)
	if err != nil {
		return err
	}
	if product.AccountID != accountId {
		return errors.New("unauthorized")
	}

	return p.repo.DeleteProduct(ctx, productId)
}
