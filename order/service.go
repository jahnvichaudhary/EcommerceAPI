package order

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/rasadov/EcommerceAPI/pkg/utils"
	"log"
	"strconv"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, totalPrice float64, products []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error)
}

type orderService struct {
	repository Repository
	producer   sarama.AsyncProducer
}

func NewOrderService(repository Repository, producer sarama.AsyncProducer) Service {
	return &orderService{repository, producer}
}

func (service orderService) Producer() sarama.AsyncProducer {
	return service.producer
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
			err = utils.SendMessageToRecommender(service, Event{
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
