package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rasadov/EcommerceMicroservices/order"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
	AccountUrl  string `envconfig:"ACCOUNT_SERVICE_URL"`
	ProductUrl  string `envconfig:"PRODUCT_SERVICE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var repository order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = order.NewPostgresRepository(cfg.DatabaseUrl)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()

	log.Println("Listening on port 8080...")
	service := order.NewOrderService(repository)
	log.Fatal(order.ListenGRPC(service, cfg.AccountUrl, cfg.ProductUrl, 8080))
}
