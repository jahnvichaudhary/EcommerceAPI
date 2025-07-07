package main

import (
	"github.com/rasadov/EcommerceAPI/payment/config"
	"github.com/rasadov/EcommerceAPI/payment/internal"
	"github.com/tinrab/retry"
	"log"
	"time"
)

func main() {
	var repository internal.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = internal.NewPostgresRepository(config.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})

	dodoClient := internal.NewDodoClient(config.DodoAPIKEY, config.DodoTestMode)

	service := internal.NewPaymentService(dodoClient, repository)

	log.Fatal(internal.StartServers(service, config.OrderServiceURL, config.GrpcPort, config.WebhookPort))
}
