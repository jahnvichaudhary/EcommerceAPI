package order

import (
	"context"
	"github.com/IBM/sarama"
	"log"
	"strconv"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, totalPrice float64, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID            uint `gorm:"primaryKey;autoIncrement"`
	CreatedAt     time.Time
	TotalPrice    float64
	AccountID     string
	productsInfos []ProductsInfo   `gorm:"foreignKey:OrderID"`
	Products      []OrderedProduct `gorm:"-"`
}

type ProductsInfo struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	OrderID   uint
	ProductID string
	Quantity  int
}

func (ProductsInfo) TableName() string {
	return "order_products"
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
	producer   sarama.AsyncProducer
}

func NewOrderService(repository Repository, producer sarama.AsyncProducer) Service {
	return &orderService{repository, producer}
}

func (service orderService) PostOrder(ctx context.Context, accountID string, totalPrice float64, products []OrderedProduct) (*Order, error) {
	order := Order{
		AccountID:  accountID,
		TotalPrice: totalPrice,
		Products:   products,
		CreatedAt:  time.Now().UTC(),
	}
	err := service.repository.PutOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// Send to recommendation service
	go func() {
		accountIdInt, err := strconv.Atoi(accountID)
		if err != nil {
			log.Println("Failed to convert account ID to int:", err)
			return
		}
		for _, product := range products {
			err = service.SendMessageToRecommender(Event{
				Type: "purchase",
				EventData: EventData{
					AccountId: accountIdInt,
					ProductId: product.ID,
				},
			}, "interaction_events")
			if err != nil {
				log.Println("Failed to send event to recommendation service:", err)
			}
		}
	}()

	return &order, nil
}

func (service orderService) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	return service.repository.GetOrdersForAccount(ctx, accountID)
}
