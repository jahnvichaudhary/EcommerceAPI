package config

import "os"

var (
	DatabaseURL     string
	DodoAPIKEY      string
	DodoCheckoutURL string
	DodoTestMode    bool
	OrderServiceURL string
)

const (
	WebhookPort int = 8081
	GrpcPort    int = 8080
)

func init() {
	DatabaseURL = os.Getenv("DATABASE_URL")
	DodoAPIKEY = os.Getenv("DODO_API_KEY")
	DodoCheckoutURL = os.Getenv("DODO_CHECKOUT_URL")
	DodoTestMode = os.Getenv("DODO_TEST_MODE") == "true"
	OrderServiceURL = os.Getenv("ORDER_SERVICE_URL")
}
