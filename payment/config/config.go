package config

import "os"

var (
	DodoAPIKEY      string
	DodoCheckoutURL string
	DodoRedirectURL string
)

func init() {
	DodoAPIKEY = os.Getenv("DODO_API_KEY")
	DodoCheckoutURL = os.Getenv("DODO_CHECKOUT_URL")
	DodoRedirectURL = os.Getenv("DODO_REDIRECT_URL")
}
