package internal

import (
	"context"
	"log"
	"net/http"
	"time"

	order "github.com/rasadov/EcommerceAPI/order/client"
)

type WebhookServer struct {
	service     Service
	orderClient *order.Client
}

type WebhookPayload struct {
	Type string `json:"type"`
	Data struct {
		Customer struct {
			CustomerID string `json:"customer_id"`
			Email      string `json:"email"`
			Name       string `json:"name"`
		} `json:"customer"`
		ProductCart []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		} `json:"product_cart"` // Product cart is going to be a slice of one element since
		// we always pass one product with the quantity one
		PaymentId string `json:"payment_id"`
	} `json:"data"`
}

func (s *WebhookServer) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	transaction, err := s.service.HandlePaymentWebhook(ctx, w, r)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = s.orderClient.UpdateOrderStatus(ctx, transaction.OrderId, transaction.Status)
	if err != nil {
		log.Println(err.Error())
	}
}
