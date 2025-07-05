package config

import "os"

var (
	AccountUrl string
	ProductUrl string
	OrderUrl   string
	SecretKey  string
	Issuer     string
)

func init() {
	AccountUrl = os.Getenv("ACCOUNT_URL")
	ProductUrl = os.Getenv("PRODUCT_URL")
	OrderUrl = os.Getenv("ORDER_URL")
	SecretKey = os.Getenv("SECRET_KEY")
	Issuer = os.Getenv("ISSUER")
}
