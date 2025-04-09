package main

import (
	"github.com/IBM/sarama"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"

	"github.com/rasadov/EcommerceAPI/product/internal"
)

type Config struct {
	DatabaseURL      string `envconfig:"DATABASE_URL"`
	BootstrapServers string `envconfig:"BOOTSTRAP_SERVERS" default:"kafka:9092"`
}

func main() {
	var cfg Config
	var repository internal.Repository
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
		repository, err = internal.NewElasticRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()
	log.Println("Listening on port 8080...")
	service := internal.NewProductService(repository, producer)
	log.Fatal(internal.ListenGRPC(service, 8080))
}
