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
	AccountUrl  string `envconfig:"ACCOUNT_URL"`
	ProductUrl  string `envconfig:"PRODUCT_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = order.NewPostgresRepository(cfg.DatabaseUrl)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := order.NewOrderService(r)
	log.Fatal(order.ListenGRPC(s, cfg.AccountUrl, cfg.ProductUrl, 8080))
}
