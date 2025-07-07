package internal

import (
	"context"
	"encoding/json"
	order "github.com/rasadov/EcommerceAPI/order/client"
	"github.com/rasadov/EcommerceAPI/payment/models"
	"io"
	"log"
	"net/http"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Verify webhook signature if your payment provider uses one

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse webhook payload
	var payload WebhookPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	var productId string
	for _, p := range payload.Data.ProductCart {
		productId = p.ProductID
	}

	// Process the webhook based on event type
	ctx := context.Background()
	switch payload.Type {
	case "payment.succeeded":
		// Save transaction details in db
		err = s.service.ProcessPayment(ctx, payload.Data.Customer.CustomerID,
			productId, payload.Data.PaymentId, models.Success)
		if err != nil {
			log.Println(err)
		}
		// TODO: Inform order microservice
	case "payment.failed":
		// Save transaction details in db
		err = s.service.ProcessPayment(ctx, payload.Data.Customer.CustomerID,
			productId, payload.Data.PaymentId, models.Failed)
		// TODO: Inform order microservice
	default:
		log.Printf("Unhandled webhook event type: %s", payload.Type)
	}

	if err != nil {
		log.Printf("Error processing webhook: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return a 200 OK to acknowledge receipt of the webhook
	w.WriteHeader(http.StatusOK)
}
