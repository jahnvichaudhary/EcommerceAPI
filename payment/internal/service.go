package internal

import (
	"context"
	"errors"
	"github.com/dodopayments/dodopayments-go"
	"github.com/rasadov/EcommerceAPI/payment/models"
	"gorm.io/gorm"
	"net/http"
)

type Service interface {
	CreateCustomerPortalSession(ctx context.Context,
		customer *models.Customer) (string, error)
	FindOrCreateCustomer(ctx context.Context,
		userId int64,
		email, name string) (*models.Customer, error)
	GetCheckoutURL(ctx context.Context,
		email, name, redirect string,
		price int64,
		currency dodopayments.Currency) (checkoutURL string, productId string, err error)
	RegisterTransaction(ctx context.Context,
		orderId, userId, price int64,
		currency dodopayments.Currency,
		customerId, productId string) error
	HandlePaymentWebhook(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.Transaction, error)
}

type paymentService struct {
	client            PaymentClient
	paymentRepository Repository
}

func NewPaymentService(client PaymentClient, paymentRepository Repository) Service {
	return &paymentService{client: client, paymentRepository: paymentRepository}
}

// GetCheckoutURL - returns url to check out page, productId and error.
// Called after creating product and registering productId with order
func (d *paymentService) GetCheckoutURL(ctx context.Context,
	email, name, redirect string,
	price int64,
	currency dodopayments.Currency) (checkoutURL string, productId string, err error) {
	return d.client.CreateCheckoutLink(ctx, email, name, redirect, price, currency)
}

func (d *paymentService) CreateCustomerPortalSession(ctx context.Context, customer *models.Customer) (string, error) {
	customerPortalLink, err := d.client.CreateCustomerSession(ctx, customer.CustomerId)
	if err != nil {
		return "", err
	}
	return customerPortalLink, nil
}

func (d *paymentService) FindOrCreateCustomer(ctx context.Context, userId int64, email, name string) (*models.Customer, error) {
	existingCustomer, err := d.paymentRepository.GetCustomerByUserID(ctx, userId)

	if err == nil {
		return existingCustomer, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	customer, err := d.client.CreateCustomer(ctx, userId, email, name)

	if err != nil {
		return nil, err
	}

	err = d.paymentRepository.SaveCustomer(ctx, customer)

	return customer, err
}

func (d *paymentService) RegisterTransaction(ctx context.Context,
	orderId, userId, price int64,
	currency dodopayments.Currency,
	customerId, productId string) error {
	transaction := &models.Transaction{
		OrderId:    orderId,
		UserId:     userId,
		CustomerId: customerId,
		ProductId:  productId,
		Amount:     price,
		Currency:   string(currency),
	}

	return d.paymentRepository.RegisterTransaction(ctx, transaction)
}

func (d *paymentService) HandlePaymentWebhook(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.Transaction, error) {
	updatedTransaction, err := d.client.HandleWebhook(w, r)
	if err != nil {
		return nil, err
	}

	transaction, err := d.paymentRepository.GetTransactionByProductID(ctx, updatedTransaction.ProductId)
	if err != nil {
		return nil, err
	}

	transaction.PaymentId = updatedTransaction.PaymentId
	transaction.Status = updatedTransaction.Status

	err = d.paymentRepository.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
