package main

import (
	"github.com/IBM/sarama"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"

	"github.com/rasadov/EcommerceAPI/product/internal/product"
	"github.com/rasadov/EcommerceAPI/product/internal/server"
)

type Config struct {
	DatabaseURL      string `envconfig:"DATABASE_URL"`
	BootstrapServers string `envconfig:"BOOTSTRAP_SERVERS" default:"kafka:9092"`
}

func main() {
	var cfg Config
	var repository product.Repository
	var producer sarama.AsyncProducer

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	producer, err = sarama.NewAsyncProducer([]string{cfg.BootstrapServers}, nil)
	if err != nil {
		log.Println(err)
	}
	defer producer.Close()

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = product.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()
	log.Println("Listening on port 8080...")
	service := product.NewProductService(repository, producer)
	log.Fatal(server.ListenGRPC(service, 8080))
}
