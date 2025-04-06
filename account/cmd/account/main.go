package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/rasadov/EcommerceAPI/account/internal/server"
	"github.com/rasadov/EcommerceAPI/account/internal/user"
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
	repository user.Repository
)

func main() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = user.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()
	jwtService := user.NewJwtService(cfg.SecretKey, cfg.Issuer)
	log.Println("Listening on port 8080...")
	service := user.NewService(repository, jwtService)
	log.Fatal(server.ListenGRPC(service, 8080))
}
