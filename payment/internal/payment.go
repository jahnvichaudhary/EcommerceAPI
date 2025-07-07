package internal

import (
	"context"
	"fmt"
	"github.com/dodopayments/dodopayments-go"
	"github.com/dodopayments/dodopayments-go/option"
	"github.com/rasadov/EcommerceAPI/payment/config"
	"github.com/rasadov/EcommerceAPI/payment/models"
)

type PaymentClient interface {
	CreateCustomer(ctx context.Context, userId int64, email, name string) (*models.Customer, error)
	CreateCheckoutLink(ctx context.Context,
		email, name, redirect string, price int64,
		currency dodopayments.Currency) (checkoutURL string, productId string, err error)
	CreateCustomerSession(ctx context.Context, customerId string) (string, error)
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
