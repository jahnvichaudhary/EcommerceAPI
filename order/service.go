package order

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type orderService struct {
	repository Repository
}

func NewOrderService(r Repository) Service {
	return &orderService{r}
}

func (o orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	order := Order{
		ID:         ksuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		AccountID:  accountID,
		TotalPrice: 0.0,
		Products:   products,
	}
	for _, product := range products {
		order.TotalPrice += product.Price
	}
	err := o.repository.PutOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return o.repository.GetOrdersForAccount(ctx, accountID)
}
