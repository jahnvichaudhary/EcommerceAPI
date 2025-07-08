package internal

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dodopayments/dodopayments-go"
	"github.com/dodopayments/dodopayments-go/option"
	"github.com/rasadov/EcommerceAPI/payment/config"
	"github.com/rasadov/EcommerceAPI/payment/models"
	"io"
	"log"
	"net/http"
)

type PaymentClient interface {
	CreateCustomer(ctx context.Context, userId int64, email, name string) (*models.Customer, error)
	CreateCheckoutLink(ctx context.Context,
		email, name, redirect string, price int64,
		currency dodopayments.Currency) (checkoutURL string, productId string, err error)
	CreateCustomerSession(ctx context.Context, customerId string) (string, error)
	HandleWebhook(w http.ResponseWriter, r *http.Request) (*models.Transaction, error)
}

func NewDodoClient(apiKey string, testMode bool) PaymentClient {
	if testMode {
		return &dodoClient{
			client: dodopayments.NewClient(
				option.WithBearerToken(apiKey),
				option.WithEnvironmentTestMode(),
			),
			webhookSecret: "",
		}
	}

	return &dodoClient{
		client: dodopayments.NewClient(
			option.WithBearerToken(apiKey),
		),
		webhookSecret: "",
	}
}

type dodoClient struct {
	client        *dodopayments.Client
	webhookSecret string
}

func (d *dodoClient) CreateCustomer(ctx context.Context, userId int64, email, name string) (*models.Customer, error) {
	customer, err := d.client.Customers.New(ctx, dodopayments.CustomerNewParams{
		Email: dodopayments.F(email),
		Name:  dodopayments.F(name),
	})

	if err != nil {
		return nil, err
	}

	return &models.Customer{
		UserId:     userId,
		CustomerId: customer.CustomerID,
		CreatedAt:  customer.CreatedAt,
	}, nil
}

func (d *dodoClient) CreateCheckoutLink(ctx context.Context,
	email, name, redirect string,
	price int64,
	currency dodopayments.Currency) (checkoutURL string, productId string, err error) {
	product, err := d.client.Products.New(ctx, dodopayments.ProductNewParams{
		Price: dodopayments.F[dodopayments.PriceUnionParam](dodopayments.PriceOneTimePriceParam{
			Currency:              dodopayments.F(currency),
			Discount:              dodopayments.F(0.000000),
			Price:                 dodopayments.F(price),
			PurchasingPowerParity: dodopayments.F(true),
			Type:                  dodopayments.F(dodopayments.PriceOneTimePriceTypeOneTimePrice),
		}),
		TaxCategory: dodopayments.F(dodopayments.TaxCategorySaas),
	})

	if err != nil {
		return "", "", err
	}

	checkoutUrl := fmt.Sprintf(
		"%s/%s?quantity=1&email=%s&disableEmail=true&fullName=%s&disableFullName=true&redirect_url=%s",
		config.DodoCheckoutURL, product.ProductID, email, name, redirect)
	return checkoutUrl, product.ProductID, nil
}

func (d *dodoClient) CreateCustomerSession(ctx context.Context, customerId string) (string, error) {
	customerPortal, err := d.client.Customers.CustomerPortal.New(ctx, customerId,
		dodopayments.CustomerCustomerPortalNewParams{})
	if err != nil {
		return "", err
	}
	return customerPortal.Link, nil
}

func (d *dodoClient) verifyWebhookSignature(signature string, payload []byte) bool {
	h := hmac.New(sha256.New, []byte(config.DodoWebhookSecret))
	h.Write(payload)
	expectedSignature := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (d *dodoClient) HandleWebhook(w http.ResponseWriter, r *http.Request) (*models.Transaction, error) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil, errors.New("method not allowed")
	}

	webhookSignature := r.Header.Get("webhook-signature")
	if !d.verifyWebhookSignature(webhookSignature, []byte(config.DodoWebhookSecret)) {
		http.Error(w, "Invalid Webhook Signature", http.StatusBadRequest)
		return nil, errors.New("invalid Webhook Signature")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return nil, err
	}
	defer r.Body.Close()

	// Parse webhook payload
	var payload WebhookPayload

	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return nil, err
	}

	var productId string
	for _, p := range payload.Data.ProductCart {
		productId = p.ProductID
	}

	transaction := &models.Transaction{
		CustomerId: payload.Data.Customer.CustomerID,
		ProductId:  productId,
		PaymentId:  payload.Data.PaymentId,
	}

	// Process the webhook based on event type
	switch payload.Type {
	case "payment.succeeded":
		transaction.Status = string(models.Success)
	case "payment.failed":
		transaction.Status = string(models.Failed)
	default:
		log.Printf("Unhandled webhook event type: %s", payload.Type)
	}

	// Return a 200 OK to acknowledge receipt of the webhook
	w.WriteHeader(http.StatusOK)
	return transaction, nil
}
