package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/rasadov/EcommerceAPI/graphql/graph"
	"github.com/rasadov/EcommerceAPI/pkg/auth"
	"github.com/rasadov/EcommerceAPI/pkg/middleware"
	"log"
)

type AppConfig struct {
	AccountUrl string `envconfig:"ACCOUNT_SERVICE_URL"`
	ProductUrl string `envconfig:"PRODUCT_SERVICE_URL"`
	OrderUrl   string `envconfig:"ORDER_SERVICE_URL"`
	SecretKey  string `envconfig:"SECRET_KEY"`
	Issuer     string `envconfig:"ISSUER"`
}

func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	server, err := graph.NewGraphQLServer(cfg.AccountUrl, cfg.ProductUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.New(server.ToExecutableSchema())
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	engine := gin.Default()

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "It works",
		})
	})
	engine.POST("/graphql", middleware.AuthorizeJWT(auth.NewJwtService(cfg.SecretKey, cfg.Issuer)), gin.WrapH(srv))
	engine.GET("/playground", gin.WrapH(playground.Handler("Playground", "/graphql")))

	log.Fatal(engine.Run(":8080"))
}
