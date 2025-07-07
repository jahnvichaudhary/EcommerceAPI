package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/dodopayments/dodopayments-go"
	"github.com/rasadov/EcommerceAPI/payment/config"
	"github.com/rasadov/EcommerceAPI/payment/models"
	"gorm.io/gorm"
)

type Service interface {
	GetCustomerPortal(ctx context.Context,
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
}

type dodoPaymentService struct {
	client            *dodopayments.Client
	paymentRepository Repository
}

func NewPaymentService(client *dodopayments.Client, paymentRepository Repository) Service {
	return &dodoPaymentService{client: client, paymentRepository: paymentRepository}
}

// GetCheckoutURL - returns url to check out page, productId and error.
// Called after creating product and registering productId with order
func (d *dodoPaymentService) GetCheckoutURL(ctx context.Context,
	email, name, redirect string,
	price int64,
	currency dodopayments.Currency) (checkoutURL string, productId string, err error) {
	productID, err := d.registerProduct(ctx, price, currency)
	if err != nil {
		return "", "", err
	}

	checkoutUrl := fmt.Sprintf("%s/%s?quantity=1&email=%s&disableEmail=true&fullName=%s&disableFullName=true&redirect_url=%s", config.DodoCheckoutURL, productID, email, name, redirect)
	return checkoutUrl, productID, nil
}

func (d *dodoPaymentService) GetCustomerPortal(ctx context.Context, customer *models.Customer) (string, error) {
	customerPortal, err := d.client.Customers.CustomerPortal.New(ctx, customer.CustomerId, dodopayments.CustomerCustomerPortalNewParams{})
	if err != nil {
		return "", err
	}
	return customerPortal.Link, nil
}

func (d *dodoPaymentService) FindOrCreateCustomer(ctx context.Context, userId int64, email, name string) (*models.Customer, error) {
	existingCustomer, err := d.paymentRepository.GetCustomerByUserID(ctx, userId)

	if err == nil {
		return existingCustomer, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	customer, err := d.client.Customers.New(ctx, dodopayments.CustomerNewParams{
		Email: dodopayments.F(email),
		Name:  dodopayments.F(name),
	})

	if err != nil {
		return nil, err
	}

	newCustomer := &models.Customer{
		UserId:     userId,
		CustomerId: customer.CustomerID,
		CreatedAt:  customer.CreatedAt,
	}

	err = d.paymentRepository.SaveCustomer(ctx, newCustomer)

	return newCustomer, err
}

// registerProduct - Used before sending checkout url, we will connect this with an order
func (d *dodoPaymentService) registerProduct(ctx context.Context, price int64, currency dodopayments.Currency) (string, error) {
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
		return "", err
	}

	return product.ProductID, err
}

func (d *dodoPaymentService) RegisterTransaction(ctx context.Context,
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
