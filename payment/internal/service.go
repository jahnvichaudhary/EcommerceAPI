package internal

import "github.com/dodopayments/dodopayments-go"

type PaymentService interface {
	GetCustomerPortal(userId string) string
	RegisterProduct(productID, title, description, category string, price int) error
	GetCheckoutURL(productID string) string
}

type dodoPaymentService struct {
	client *dodopayments.Client
}

func NewPaymentService(client *dodopayments.Client) PaymentService {
	return &dodoPaymentService{client: client}
}

func (d *dodoPaymentService) GetCheckoutURL(productID string) string {
	return ""
}

func (d *dodoPaymentService) GetCustomerPortal(userId string) string {
	return ""
}

func (d *dodoPaymentService) RegisterProduct(productID, title, description, category string, price int) error {
	return nil
}
