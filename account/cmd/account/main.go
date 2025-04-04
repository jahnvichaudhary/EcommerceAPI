package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rasadov/EcommerceAPI/account"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL" default:"postgres://user:password@localhost/dbname?sslmode=disable"`
	SecretKey   string `envconfig:"SECRET_KEY"`
	Issuer      string `envconfig:"ISSUER"`
}

var (
	cfg        Config
	repository account.Repository
)

func main() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = account.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()
	jwtService := account.NewJwtService(cfg.SecretKey, cfg.Issuer)
	log.Println("Listening on port 8080...")
	service := account.NewService(repository, jwtService)
	log.Fatal(account.ListenGRPC(service, 8080))
}
