package models

type Customer struct {
	CustomerId string `json:"customer_id"`
	UserId     string `json:"user_id"`

	Transactions []Transaction `json:"transactions"`
}

type CustomerInput struct {
	UserId string `json:"user_id"`
}
