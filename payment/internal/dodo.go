package internal

import (
	"github.com/dodopayments/dodopayments-go"
	"github.com/dodopayments/dodopayments-go/option"
)

func NewDodoClient(apiKey string, testMode bool) *dodopayments.Client {
	if testMode {
		return dodopayments.NewClient(
			option.WithBearerToken(apiKey),
			option.WithEnvironmentTestMode(),
		)
	}

	return dodopayments.NewClient(
		option.WithBearerToken(apiKey),
	)
}
