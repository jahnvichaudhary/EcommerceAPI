package product

import (
	"context"
	"encoding/json"
	"errors"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

var (
	ErrNotFound = errors.New("entity not found")
)

type Repository interface {
	Close()
	PutProduct(ctx context.Context, p Product) error
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
	UpdateProduct(ctx context.Context, updatedProduct Product) error
	DeleteProduct(ctx context.Context, productId string) error
}

type elasticRepository struct {
	client *elastic.Client
}

type productDocument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

func NewElasticRepository(url string) (Repository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &elasticRepository{client}, nil
}

func (r *elasticRepository) Close() {
	r.client.Stop()
}

func (r *elasticRepository) PutProduct(ctx context.Context, p Product) error {
	_, err := r.client.Index().
		Index("catalog").
		Type("product").
		BodyJson(productDocument{
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		}).
		Do(ctx)
	return err
}

func (r *elasticRepository) GetProductById(ctx context.Context, id string) (*Product, error) {
	res, _ := r.client.Get().
		Index("catalog").
		Type("product").
		Id(id).
		Do(ctx)
	if !res.Found {
		return nil, ErrNotFound
	}
	product := productDocument{}
	if err := json.Unmarshal(*res.Source, &product); err != nil {
		return nil, err
	}
	return &Product{
		ID:          id,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

func (r *elasticRepository) ListProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.MatchAllQuery{}).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []Product
	for _, hit := range res.Hits.Hits {
		product := productDocument{}
		if err = json.Unmarshal(*hit.Source, &product); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, err
}

func (r *elasticRepository) ListProductsWithIDs(ctx context.Context, ids []string) ([]Product, error) {
	var items []*elastic.MultiGetItem
	for _, id := range ids {
		items = append(items, elastic.NewMultiGetItem().
			Index("catalog").
			Type("product").
			Id(id))
	}
	res, err := r.client.MultiGet().
		Add(items...).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []Product
	for _, doc := range res.Docs {
		product := productDocument{}
		if err = json.Unmarshal(*doc.Source, &product); err == nil {
			products = append(products, Product{
				ID:          doc.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, err
}

func (r *elasticRepository) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	res, err := r.client.Search().
		Index("catalog").
		Type("product").
		Query(elastic.NewMultiMatchQuery(query, "name", "description")).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var products []Product
	for _, hit := range res.Hits.Hits {
		product := productDocument{}
		if err = json.Unmarshal(*hit.Source, &product); err == nil {
			products = append(products, Product{
				ID:          hit.Id,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
			})
		}
	}
	return products, err
}

func (r *elasticRepository) UpdateProduct(ctx context.Context, updatedProduct Product) error {
	_, err := r.client.Update().
		Index("catalog").
		Type("product").
		Id(updatedProduct.ID).
		Doc(productDocument{
			Name:        updatedProduct.Name,
			Description: updatedProduct.Description,
			Price:       updatedProduct.Price,
		}).
		Do(ctx)
	return err
}

func (r *elasticRepository) DeleteProduct(ctx context.Context, productId string) error {
	_, err := r.client.Delete().
		Index("catalog").
		Type("product").
		Id(productId).
		Do(ctx)
	return err
}
