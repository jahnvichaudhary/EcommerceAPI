package main

import (
	"github.com/rasadov/EcommerceAPI/account/config"
	"github.com/rasadov/EcommerceAPI/account/internal"
	"github.com/rasadov/EcommerceAPI/pkg/auth"
	"github.com/tinrab/retry"
	"log"
	"time"
)

var (
	repository internal.Repository
)

func main() {
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repository, err = internal.NewPostgresRepository(config.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer repository.Close()
	jwtService := auth.NewJwtService(config.SecretKey, config.Issuer)
	log.Println("Listening on port 8080...")
	service := internal.NewService(repository, jwtService)
	log.Fatal(internal.ListenGRPC(service, 8080))
}
