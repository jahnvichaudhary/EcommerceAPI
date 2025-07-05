package models

import "time"

type Transaction struct {
	PaymentID  string    `json:"id" gorm:"primary_key"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;"`
	CustomerId string    `json:"customer_id"`
	OrderId    string    `json:"order_id"`
	Amount     int       `json:"amount"`
	Currency   string    `json:"currency"`
}
