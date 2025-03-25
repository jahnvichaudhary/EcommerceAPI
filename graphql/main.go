package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type AppConfig struct {
	AccountUrl string `envconfig:"ACCOUNT_SERVICE_URL"`
	ProductUrl string `envconfig:"PRODUCT_SERVICE_URL"`
	OrderUrl   string `envconfig:"ORDER_SERVICE_URL"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(cfg.AccountUrl, cfg.ProductUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.New(s.toExecutableSchema())
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("Rauf", "/graphql"))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
