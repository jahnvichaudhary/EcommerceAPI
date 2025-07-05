package config

import "os"

var (
	DatabaseUrl      string
	AccountUrl       string
	ProductUrl       string
	BootstrapServers string
)

func init() {
	DatabaseUrl = os.Getenv("DATABASE_URL")
	AccountUrl = os.Getenv("ACCOUNT_URL")
	ProductUrl = os.Getenv("PRODUCT_URL")
	BootstrapServers = os.Getenv("BOOTSTRAP_SERVERS")
}
