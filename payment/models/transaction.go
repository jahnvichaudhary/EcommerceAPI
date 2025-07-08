package models

import (
	"time"
)

type TransactionStatus string

const (
	Failed  = TransactionStatus("Failed")
	Success = TransactionStatus("Success")
)

func (s TransactionStatus) String() string {
	return string(s)
}

type Transaction struct {
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;"`
	OrderId    int64     `json:"order_id"`
	UserId     int64     `json:"user_id"`
	CustomerId string    `json:"customer_id"`
	ProductId  string    `json:"product_id" gorm:"primaryKey;"`
	PaymentId  string    `json:"payment_id"`
	Amount     int64     `json:"amount"`
	Currency   string    `json:"currency"`
	Status     string    `json:"status" gorm:"type:varchar(20)"`
}
