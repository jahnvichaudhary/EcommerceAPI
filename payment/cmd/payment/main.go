package main

import (
	"github.com/rasadov/EcommerceAPI/payment/config"
	"github.com/rasadov/EcommerceAPI/payment/internal"
	"github.com/tinrab/retry"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func main() {
	var repository internal.Repository

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
		if err != nil {
			log.Println(err)
		}
		repository, err = internal.NewPostgresRepository(db)
		if err != nil {
			log.Println(err)
		}
		return
	})

	dodoClient := internal.NewDodoClient(config.DodoAPIKEY, config.DodoTestMode)

	service := internal.NewPaymentService(dodoClient, repository)

	log.Fatal(internal.StartServers(service, config.OrderServiceURL, config.GrpcPort, config.WebhookPort))
}
