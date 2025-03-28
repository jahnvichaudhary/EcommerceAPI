package order

import (
	"context"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}

type Product struct {
	OrderID   uint
	ProductID uint
	Quantity  int
}

func (Product) TableName() string {
	return "order_products"
}

type OrderedProduct struct {
	ID          uint
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
		AccountID:  accountID,
		TotalPrice: 0.0,
		Products:   products,
		CreatedAt:  time.Now().UTC(),
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
