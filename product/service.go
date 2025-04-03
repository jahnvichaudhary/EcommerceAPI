package product

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"log"
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
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	AccountID   int     `json:"accountID"`
}

type productService struct {
	repo     Repository
	producer sarama.AsyncProducer
}

func NewProductService(repository Repository, producer sarama.AsyncProducer) Service {
	return &productService{repository, producer}
}

func (service productService) PostProduct(ctx context.Context, name, description string, price float64, accountId int) (*Product, error) {
	product := Product{
		Name:        name,
		Description: description,
		Price:       price,
		AccountID:   accountId,
	}

	err := service.repo.PutProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	go func() {
		err = service.SendMessageToRecommender(Event{
			Type: "product_created",
			Data: EventData{
				ID:          &product.ID,
				Name:        &product.Name,
				Description: &product.Description,
				Price:       &product.Price,
				AccountID:   &product.AccountID,
			},
		}, "product_events")
		if err != nil {
			log.Println("Failed to send event to recommendation service:", err)
		}
	}()

	return &product, nil
}

func (service productService) GetProduct(ctx context.Context, id string) (*Product, error) {
	product, err := service.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}

	go func() {
		err = service.SendMessageToRecommender(Event{
			Type: "product_retrieved",
			Data: EventData{
				ID:        &product.ID,
				AccountID: &product.AccountID,
			},
		}, "interaction_events")
		if err != nil {
			log.Println("Failed to send event to recommendation service:", err)
		}
	}()

	return product, nil
}

func (service productService) GetProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	return service.repo.ListProducts(ctx, skip, take)
}

func (service productService) GetProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	return service.repo.ListProductsWithIDs(ctx, ids)
}

func (service productService) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	return service.repo.SearchProducts(ctx, query, skip, take)
}

func (service productService) UpdateProduct(ctx context.Context, id, name, description string, price float64, accountId int) (*Product, error) {
	product, err := service.repo.GetProductById(ctx, id)
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
		price,
		accountId,
	}
	err = service.repo.UpdateProduct(ctx, updatedProduct)
	if err != nil {
		return nil, err
	}

	go func() {
		err = service.SendMessageToRecommender(Event{
			Type: "product_updated",
			Data: EventData{
				ID:          &updatedProduct.ID,
				Name:        &updatedProduct.Name,
				Description: &updatedProduct.Description,
				Price:       &updatedProduct.Price,
				AccountID:   &updatedProduct.AccountID,
			},
		}, "product_events")
		if err != nil {
			log.Println("Failed to send event to recommendation service:", err)
		}
	}()

	return &updatedProduct, nil
}
func (service productService) DeleteProduct(ctx context.Context, productId string, accountId int) error {
	product, err := service.repo.GetProductById(ctx, productId)
	if err != nil {
		return err
	}
	if product.AccountID != accountId {
		return errors.New("unauthorized")
	}

	go func() {
		err = service.SendMessageToRecommender(Event{
			Type: "product_deleted",
			Data: EventData{
				ID: &product.ID,
			},
		}, "product_events")
		if err != nil {
			log.Println("Failed to send event to recommendation service:", err)
		}
	}()

	return service.repo.DeleteProduct(ctx, productId)
}
